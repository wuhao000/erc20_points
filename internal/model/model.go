package model

type User struct {
  Id   uint64 `gorm:"primary_key"`
  Name string
}

type PointRecord struct {
  Id          uint64 `gorm:"primary_key"`
  NetworkId   uint64
  Address     string
  UserId      uint64
  Balance     uint64
  Points      float64
  LastUpdated uint64
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
  Time        uint64
  FromAddress string `gorm:"index:idx_from_address_address"`
  ToAddress   string `gorm:"index:idx_to_address_address"`
  Value       uint64
  CalcStatus  uint8
}

type Event struct {
  Id          uint64 `gorm:"primary_key"`
  Hash        string `gorm:"unique"`
  Time        uint64
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
  BalanceOrigin uint64
  BalanceNew    uint64
  PointsNew     float64
  StartTime     uint64
  EndTime       uint64
}
