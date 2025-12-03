package service

import (
  "fmt"
  "testing"

  "com.wh1200.points/internal/config"
  "com.wh1200.points/internal/model"
)

func TestCalc(t *testing.T) {
  zeroAddress := "0x0"
  address := "0x123"
  anotherAddress := "0x456"
  transferRecords := make([]model.TransferRecord, 0)

  transferRecords = append(transferRecords, model.TransferRecord{
    FromAddress: zeroAddress,
    ToAddress:   address,
    Value:       50,
    BlockNumber: 100,
    Time:        1000,
  })

  transferRecords = append(transferRecords, model.TransferRecord{
    FromAddress: zeroAddress,
    ToAddress:   anotherAddress,
    Value:       50,
    BlockNumber: 101,
    Time:        2000,
  })

  transferRecords = append(transferRecords, model.TransferRecord{
    FromAddress: anotherAddress,
    ToAddress:   address,
    Value:       50,
    BlockNumber: 102,
    Time:        3000,
  })

  cfg := &config.ChainConfig{}
  cfg.ChainId = 1
  pr := &model.PointRecord{}
  pr.Balance = 0
  pr.NetworkId = 1
  pr.Points = 0
  a, b := calc(cfg, address, pr, transferRecords)
  for _, i := range a {
    fmt.Printf("%+v\n", i)
  }
  fmt.Printf("%+v\n", b)
}
