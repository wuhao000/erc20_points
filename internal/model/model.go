package model

import "time"

type User struct {
  Id   uint64 `gorm:"primary_key"`
  Name string
}

type PointRecord struct {
  Id          uint64 `gorm:"primary_key"`
  NetworkId   uint64
  Address     string
  UserId      uint64
  Balance     int64
  Points      float64
  LastUpdated time.Time
}

type Network struct {
  Id   uint64 `gorm:"primary_key"`
  Name string
}

type TransferRecord struct {
  Id          uint64 `gorm:"primary_key"`
  Hash        string `gorm:"unique"`
  NetworkId   uint64
  BlockNumber uint64
  ChainId     uint64
  Time        time.Time
  FromAddress string `gorm:"index:idx_from_address_address"`
  ToAddress   string `gorm:"index:idx_to_address_address"`
  Value       int64
  CalcStatus  uint8
}

type Event struct {
  Id          uint64 `gorm:"primary_key"`
  Hash        string `gorm:"unique"`
  Time        time.Time
  ChainId     uint64
  BlockNumber uint64
  NetworkId   uint64
  Name        string
  Data        string
}

type PointChangeRecord struct {
  Id            uint64 `gorm:"primary_key"`
  NetworkId     uint64
  Address       string
  PointsOrigin  float64
  BalanceOrigin int64
  BalanceNew    int64
  PointsNew     float64
  StartTime     time.Time
  EndTime       time.Time
  TransferHash  string
}
