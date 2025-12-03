package service

import (
  "context"
  "encoding/json"
  "fmt"
  "math/big"
  "strconv"
  "time"

  "com.wh1200.points/internal/chain"
  "com.wh1200.points/internal/config"
  "com.wh1200.points/internal/model"
  "github.com/ethereum/go-ethereum"
  "github.com/ethereum/go-ethereum/common"
  "github.com/ethereum/go-ethereum/core/types"
  "github.com/ethereum/go-ethereum/crypto"
  "github.com/zeromicro/go-zero/core/threading"
)

const (
  APPROVAL_HASH = "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925"
  TRANSFER_HASH = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
  MINT_HASH     = "0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885"
  BURN_HASH     = "0xcc16f5dbb4873280815c1ee09dbd06736cffcc184412cf7a71a0fdb75d397ca5"
)

func Start() {
  for _, chainConfig := range config.Cfg.Chains {
    startSyncTask(&chainConfig)
    startCalcPointsTask(&chainConfig)
  }
}

func startSyncTask(config *config.ChainConfig) {
  kv := getKvStore()
  threading.GoSafe(func() {
    interval := config.PollingInterval

    if interval <= 0 {
      interval = 10
    }

    sem := make(chan struct{}, 1) // 最大同时运行1个任务
    for {
      sem <- struct{}{} // 占用
      threading.GoSafe(func() {
        defer func() { <-sem }() // 释放
        defer func() {
          if r := recover(); r != nil {
            fmt.Printf("[%v] 同步任务异常: %v\n", config.Name, r)
          }
        }()
        lastBlockStr, err := kv.Get(config.Name + ":last_block")
        if err != nil {
          panic(err)
        }
        startBlock := config.StartBlock
        if lastBlockStr != "" {
          s, err := strconv.Atoi(lastBlockStr)
          if err != nil {
            panic(err)
          }
          startBlock = uint64(s) + 1
        }
        client, err := chain.New(config.ChainId, config.NodeUrl2)
        if err != nil {
          panic(err)
        }
        latestBlockNum, err := client.Client.BlockNumber(context.Background())
        if err != nil {
          panic(err)
        }
        endBlock := latestBlockNum - config.DelayBlocks
        if startBlock > endBlock {
          return
        }
        fmt.Printf("任务开始: %v - %v\n", startBlock, endBlock)
        tmpEnd := startBlock
        var blocksPerQuery uint64 = 500
        allEvents := make([]model.Event, 0)
        allRecords := make([]model.TransferRecord, 0)
        for ; tmpEnd <= endBlock; {
          if tmpEnd > endBlock {
            tmpEnd = endBlock
          }
          events, transferRecords := syncBlock(client, tmpEnd, tmpEnd+blocksPerQuery, config)
          allRecords = append(allRecords, transferRecords...)
          allEvents = append(allEvents, events...)
          tmpEnd += blocksPerQuery
        }
        eventRepository.SaveAll(allEvents)
        transferRecordRepository.SaveAll(allRecords)
        err = kv.Set(config.Name+":last_block", strconv.Itoa(int(endBlock)))
        block, err := client.Client.BlockByNumber(context.Background(), big.NewInt(int64(endBlock)))
        err = kv.Set(config.Name+":last_block_time", strconv.Itoa(int(block.Time())))
        if err != nil {
          panic(err)
        }
        fmt.Printf("[%s]任务结束, 从[%v]到[%v]! \n", config.Name, startBlock, endBlock)
      })
      time.Sleep(time.Duration(interval) * time.Second) // 固定等待
    }
  })
}

func syncBlock(
  client *chain.Client,
  fromBlock,
  toBlock uint64,
  chainConfig *config.ChainConfig,
) ([]model.Event, []model.TransferRecord) {
  contract := common.HexToAddress(chainConfig.ContractAddress)
  query := ethereum.FilterQuery{
    Addresses: []common.Address{contract},
  }
  query.FromBlock = big.NewInt(int64(fromBlock))
  query.ToBlock = big.NewInt(int64(toBlock))
  logs, err := client.Client.FilterLogs(context.Background(), query)
  events := make([]model.Event, 0)
  if err != nil {
    panic(err)
  }
  transferSeen := make(map[string]bool)
  eventSeen := make(map[string]bool)
  transferRecords := make([]model.TransferRecord, 0)
  for _, log := range logs {
    eventHash := log.Topics[0].Hex()
    modelEvent := model.Event{}
    modelEvent.BlockNumber = log.BlockNumber
    modelEvent.ChainId = chainConfig.ChainId
    modelEvent.NetworkId = chainConfig.ChainId
    modelEvent.Hash = logUniqueID(log).String()
    modelEvent.Time = log.BlockTimestamp
    var evt interface{}

    switch eventHash {
    case APPROVAL_HASH:
      modelEvent.Name = "Approval"
      evt = parseApprovalEvent(log)
    case TRANSFER_HASH:
      modelEvent.Name = "Transfer"
      e := parseTransferEvent(log)
      evt = e
      transferRecord := model.TransferRecord{}
      transferRecord.BlockNumber = log.BlockNumber
      transferRecord.NetworkId = chainConfig.ChainId
      transferRecord.ChainId = chainConfig.ChainId
      transferRecord.FromAddress = e.From.String()
      transferRecord.ToAddress = e.To.String()
      transferRecord.Value = e.Value.Uint64()
      transferRecord.Time = log.BlockTimestamp
      transferRecord.Hash = modelEvent.Hash
      if _, ok := transferSeen[transferRecord.Hash]; !ok {
        transferRecords = append(transferRecords, transferRecord)
        transferSeen[transferRecord.Hash] = true
      }
    case MINT_HASH:
      modelEvent.Name = "Mint"
      evt = parseMintEvent(log)
    case BURN_HASH:
      modelEvent.Name = "Burn"
      evt = parseBurnEvent(log)
    }
    b, _ := json.Marshal(evt)
    modelEvent.Data = string(b)
    if _, ok := eventSeen[modelEvent.Hash]; !ok {
      events = append(events, modelEvent)
      eventSeen[modelEvent.Hash] = true
    }
  }
  return events, transferRecords
}

func parseTransferEvent(log types.Log) *model.TransferEvent {
  event := &model.TransferEvent{}
  // indexed 参数: Topics[1..]
  event.From = common.HexToAddress(log.Topics[1].Hex())
  event.To = common.HexToAddress(log.Topics[2].Hex())
  event.Value = new(big.Int).SetBytes(log.Data)
  return event
}

func parseApprovalEvent(log types.Log) *model.ApprovalEvent {
  event := &model.ApprovalEvent{}
  // indexed 参数: Topics[1..]
  event.Owner = common.HexToAddress(log.Topics[1].Hex())
  event.Spender = common.HexToAddress(log.Topics[2].Hex())
  event.Value = new(big.Int).SetBytes(log.Data)
  return event
}

func parseMintEvent(log types.Log) *model.MintEvent {
  event := &model.MintEvent{}
  event.Addr = common.BytesToAddress(log.Data[12:32])
  event.Amount = new(big.Int).SetBytes(log.Data[32:])
  return event
}

func parseBurnEvent(log types.Log) *model.BurnEvent {
  event := &model.BurnEvent{}
  event.Addr = common.BytesToAddress(log.Data[12:32])
  event.Amount = new(big.Int).SetBytes(log.Data[32:])
  return event
}

func logUniqueID(log types.Log) common.Hash {
  idx := big.NewInt(int64(log.Index)).Bytes()
  data := append(log.BlockHash.Bytes(), log.TxHash.Bytes()...)
  data = append(data, idx...)
  return crypto.Keccak256Hash(data)
}
