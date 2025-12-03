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
)

func startCalcPointsTask(cfg *config.ChainConfig) {
  c := cron.New(cron.WithSeconds())
  _, err := c.AddFunc("0 0/5 * * * ?", func() {
    execute(cfg)
  })
  if err != nil {
    panic(err)
  }
  c.Start()
}

func execute(cfg *config.ChainConfig) {
  log.Info("积分计算任务开始：" + time.Now().String())
  records := transferRecordRepository.GetAllNotCalc(cfg.ChainId)
  if len(records) == 0 {
    return
  }
  addresses := getAddresses(records)
  pointRecords := pointRepository.FindByNetworkIdAndAddresses(cfg.ChainId, addresses)
  pointMap := make(map[string]*model.PointRecord)
  for _, r := range pointRecords {
    pointMap[r.Address] = &r
  }
  allChangeRecords := make([]model.PointChangeRecord, 0)
  allPointRecords := make([]*model.PointRecord, 0)

  for _, address := range addresses {
    changeRecords, pointRecord := calc(cfg, address, pointMap[address], filterByAddress(address, records))
    allChangeRecords = append(allChangeRecords, changeRecords...)
    allPointRecords = append(allPointRecords, pointRecord)
  }
  err := save(allPointRecords, allChangeRecords, records)
  if err != nil {
    panic(err)
  }
}

func save(prs []*model.PointRecord, pcrs []model.PointChangeRecord, transferRecords []model.TransferRecord) error {
  tx := db.Begin()
  for _, t := range transferRecords {
    t.CalcStatus = 1
  }
  repository.NewTransferRecordRepository(tx).SaveAll(transferRecords)
  repository.NewPointRecordRepository(tx).SaveOrUpdateAll(prs)
  repository.NewPointChangeRecordRepository(tx).SaveAll(pcrs)
  return tx.Commit().Error
}

func calc(cfg *config.ChainConfig, address string, pr *model.PointRecord, transferRecords []model.TransferRecord) ([]model.PointChangeRecord, *model.PointRecord) {
  changeRecords := make([]model.PointChangeRecord, 0)
  if pr == nil {
    pr = &model.PointRecord{}
    pr.Address = address
    pr.NetworkId = cfg.ChainId
  }
  for _, t := range transferRecords {
    start := pr.Balance
    startTime := pr.LastUpdated
    if t.FromAddress == address {
      // 转出
      pr.Balance -= t.Value
    } else if t.ToAddress == address {
      // 转入
      pr.Balance += t.Value
    } else {
      continue
    }
    pr.LastUpdated = t.Time
    if start > 0 {
      cr := calcSingle(pr, start, startTime, t.Time)
      changeRecords = append(changeRecords, cr)
      printCr(cr)
    }
  }
  maxTimestamp := getMaxTimestamp(cfg)
  if maxTimestamp > pr.LastUpdated {
    now := uint64(time.Now().Unix())
    cr := calcSingle(pr, pr.Balance, pr.LastUpdated, now)
    pr.LastUpdated = now
    changeRecords = append(changeRecords, cr)
    printCr(cr)
  }
  return changeRecords, pr
}

func getMaxTimestamp(cfg *config.ChainConfig) uint64 {
  kv := getKvStore()
  timeStr, err := kv.Get(cfg.Name + ":last_block_time")
  if err == nil {
    panic(err)
  }
  if timeStr == "" {
    return 0
  }
  lastBlockTime, err := strconv.Atoi(timeStr)
  if err != nil {
    panic(err)
  }
  return uint64(lastBlockTime)
}

func printCr(cr model.PointChangeRecord) {
  fmt.Printf("地址: %s, 余额从 %v - %v, 积分从： %f - %f, 时间：%v min \n", cr.Address, cr.BalanceOrigin, cr.BalanceNew,
    cr.PointsOrigin, cr.PointsNew, (cr.EndTime-cr.StartTime)/60)
}

func calcSingle(pr *model.PointRecord, startBalance uint64, startTime, endTime uint64) model.PointChangeRecord {
  oldPoints := pr.Points
  newPoints := pr.Points + float64(startBalance)*0.05*float64(endTime-startTime)/60.0/60.0
  pr.Points += newPoints

  cr := model.PointChangeRecord{
    NetworkId:     pr.NetworkId,
    Address:       pr.Address,
    PointsOrigin:  oldPoints,
    PointsNew:     newPoints,
    BalanceOrigin: startBalance,
    BalanceNew:    pr.Balance,
    StartTime:     startTime,
    EndTime:       endTime,
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
    if !seen[r.FromAddress] {
      addresses = append(addresses, r.FromAddress)
      seen[r.FromAddress] = true
    }
    if !seen[r.ToAddress] {
      addresses = append(addresses, r.ToAddress)
      seen[r.ToAddress] = true
    }
  }
  return addresses
}
