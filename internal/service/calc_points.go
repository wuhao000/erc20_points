package service

import (
  "fmt"
  "strconv"
  "time"

  "com.wh1200.points/internal/config"
  "com.wh1200.points/internal/model"
  "com.wh1200.points/internal/repository"
  "github.com/ethereum/go-ethereum/log"
  "github.com/robfig/cron/v3"
  "github.com/samber/lo"
)

const (
  // ZeroAddress 0地址，计算积分要排除，没有意义
  ZeroAddress = "0x0000000000000000000000000000000000000000"
)

// 按网络创建积分计算的定时任务
func startCalcPointsTask(cfg *config.ChainConfig) {
  c := cron.New(cron.WithSeconds())
  _, err := c.AddFunc("0/10 * * * * ?", func() {
    execute(cfg)
  })
  if err != nil {
    panic(err)
  }
  c.Start()
}

// 执行积分计算
// 获取未计算到转账记录，获取当前网络的所有积分（含余额）记录
// 用转账记录计算余额变化，并计算积分
// 按已同步的最后区块时间计算积分
func execute(cfg *config.ChainConfig) {
  log.Info("积分计算任务开始：" + time.Now().String())
  records := transferRecordRepository.GetAllNotCalc(cfg.ChainId)
  allChangeRecords := make([]model.PointChangeRecord, 0)
  allPointRecords := make([]*model.PointRecord, 0)
  addresses := getAddresses(records)
  pointRecordRepository.FindAddressesByNetworkId(cfg.ChainId)

  pointRecords := pointRecordRepository.FindByNetworkId(cfg.ChainId)

  // 如果是新账户，之前没有转账记录，也就没有积分记录，还有一种情况是有积分记录，但是本次没有转账记录，
  // 所以要处理的地址列表是二者合并后的，没有转账记录的也要按时间累积积分
  addresses = append(addresses, lo.Map(pointRecords, func(item model.PointRecord, index int) string {
    return item.Address
  })...)

  pointMap := createPointRecordMap(pointRecords)
  // 按地址处理
  for _, address := range addresses {
    changeRecords, pointRecord := calc(cfg, address, pointMap[address], filterByAddress(address, records))
    allChangeRecords = append(allChangeRecords, changeRecords...)
    allPointRecords = append(allPointRecords, pointRecord)
  }
  // 保存结果，包括积分变化记录，总积分（余额）记录，以及已处理转账记录的状态
  err := save(allPointRecords, allChangeRecords, records)
  if err != nil {
    panic(err)
  }
}

func createPointRecordMap(pointRecords []model.PointRecord) map[string]*model.PointRecord {
  pointMap := make(map[string]*model.PointRecord)
  for _, r := range pointRecords {
    pointMap[r.Address] = &r
  }
  return pointMap
}

func save(prs []*model.PointRecord, pcrs []model.PointChangeRecord, transferRecords []model.TransferRecord) error {
  tx := db.Begin()
  ids := lo.Map(transferRecords, func(item model.TransferRecord, index int) uint64 {
    return item.Id
  })
  repository.NewTransferRecordRepository(tx).UpdateCalcStatusByIds(ids, 1)
  repository.NewPointRecordRepository(tx).SaveOrUpdateAll(prs)
  repository.NewPointChangeRecordRepository(tx).SaveAll(pcrs)
  return tx.Commit().Error
}

func calc(cfg *config.ChainConfig, address string, pr *model.PointRecord, transferRecords []model.TransferRecord) ([]model.PointChangeRecord, *model.PointRecord) {
  changeRecords := make([]model.PointChangeRecord, 0)
  if pr == nil {
    pr = &model.PointRecord{
      NetworkId:   cfg.ChainId,
      Address:     address,
      Balance:     0,
      Points:      0,
      LastUpdated: time.Unix(0, 0),
    }
  }
  // 按照每一次转账记录计算余额变化，并更新积分
  for _, t := range transferRecords {
    startBalance := pr.Balance
    startTime := pr.LastUpdated
    endBalance := pr.Balance
    if t.FromAddress == address {
      // 转出
      endBalance -= t.Value
    } else if t.ToAddress == address {
      // 转入
      endBalance += t.Value
    } else {
      continue
    }
    pr.Balance = endBalance
    pr.LastUpdated = t.Time
    cr := calcSingle(pr, startBalance, startTime, t.Time, t.Hash)
    changeRecords = append(changeRecords, cr)
    printCr(cr, t.Value)
  }

  // 按时间累积的积分计算
  maxTimestamp := getMaxTimestamp(cfg)
  if maxTimestamp.Unix() > pr.LastUpdated.Unix() {
    startTime := pr.LastUpdated
    endTime := maxTimestamp
    cr := calcSingle(pr, pr.Balance, startTime, endTime, "")
    pr.LastUpdated = endTime
    changeRecords = append(changeRecords, cr)
    printCr(cr, 0)
  }
  return changeRecords, pr
}

// 从缓存获取已同步的最大区块的时间戳
func getMaxTimestamp(cfg *config.ChainConfig) time.Time {
  kv := getKvStore()
  timeStr, err := kv.Get(cfg.Name + ":last_block_time")
  if err != nil {
    panic(err)
  }
  if timeStr == "" {
    return time.Unix(0, 0)
  }
  lastBlockTime, err := strconv.Atoi(timeStr)
  if err != nil {
    panic(err)
  }
  return time.Unix(int64(lastBlockTime), 0)
}

func printCr(cr model.PointChangeRecord, value int64) {
  fmt.Printf("地址: %s, 变化金额: %d 余额从 %d - %d, 积分从： %f - %f, 时间：%v min \n", cr.Address, value,
    cr.BalanceOrigin,
    cr.BalanceNew,
    cr.PointsOrigin, cr.PointsNew, (cr.EndTime.Unix()-cr.StartTime.Unix())/60)
}

func calcSingle(pr *model.PointRecord, startBalance int64, startTime, endTime time.Time, hash string) model.PointChangeRecord {
  oldPoints := pr.Points
  newPoints := pr.Points + float64(startBalance)*0.05*float64(endTime.Unix()-startTime.Unix())/60.0/60.0
  pr.Points = newPoints

  cr := model.PointChangeRecord{
    NetworkId:     pr.NetworkId,
    Address:       pr.Address,
    PointsOrigin:  oldPoints,
    PointsNew:     newPoints,
    BalanceOrigin: startBalance,
    BalanceNew:    pr.Balance,
    StartTime:     startTime,
    EndTime:       endTime,
    TransferHash:  hash,
  }

  return cr
}

func filterByAddress(address string, records []model.TransferRecord) []model.TransferRecord {
  res := make([]model.TransferRecord, 0)
  for _, r := range records {
    if r.FromAddress == address || r.ToAddress == address {
      res = append(res, r)
    }
  }
  return res
}

func getAddresses(records []model.TransferRecord) []string {
  addresses := make([]string, 0)
  seen := make(map[string]bool)
  for _, r := range records {
    if !seen[r.FromAddress] && r.FromAddress != ZeroAddress {
      addresses = append(addresses, r.FromAddress)
      seen[r.FromAddress] = true
    }
    if !seen[r.ToAddress] && r.ToAddress != ZeroAddress {
      addresses = append(addresses, r.ToAddress)
      seen[r.ToAddress] = true
    }
  }
  return addresses
}
