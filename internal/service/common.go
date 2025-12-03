package service

import (
  "com.wh1200.points/internal/repository"
  "gorm.io/gorm"
)

var db *gorm.DB

var userRepository repository.UserRepository

var pointRepository repository.PointRecordRepository

var eventRepository repository.EventRepository

var networkRepository repository.NetworkRepository

var pointChangeRecordRepository repository.PointChangeRecordRepository

var transferRecordRepository repository.TransferRecordRepository

func InjectRepositories(_db *gorm.DB) {
  db = _db
  userRepository = repository.NewUserRepository(_db)
  pointRepository = repository.NewPointRecordRepository(_db)
  eventRepository = repository.NewEventRepository(_db)
  networkRepository = repository.NewNetworkRepository(_db)
  pointChangeRecordRepository = repository.NewPointChangeRecordRepository(_db)
  transferRecordRepository = repository.NewTransferRecordRepository(_db)
}
