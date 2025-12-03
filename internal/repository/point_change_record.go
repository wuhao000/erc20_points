package repository

import (
  "com.wh1200.points/internal/model"
  "gorm.io/gorm"
)

type PointChangeRecordRepository interface {
  Create(record *model.PointChangeRecord)
  GetById(id uint64) (*model.PointChangeRecord, error)
  Update(record *model.PointChangeRecord)
  Delete(id uint64) error
  GetByNetworkId(networkId uint64) ([]model.PointChangeRecord, error)
  GetByBlockNumber(blockNumber uint64) ([]model.PointChangeRecord, error)
  GetByAddress(address string) ([]model.PointChangeRecord, error)
  GetByNetworkAndBlock(networkId, blockNumber uint64) ([]model.PointChangeRecord, error)
  GetByNetworkAndAddress(networkId uint64, address string) ([]model.PointChangeRecord, error)
  SaveAll(pcrs []model.PointChangeRecord)
}

type PointChangeRecordRepositoryImpl struct {
  db *gorm.DB
}

func (p PointChangeRecordRepositoryImpl) SaveAll(pcrs []model.PointChangeRecord) {
  if len(pcrs) == 0 {
    return
  }
  p.db.Create(&pcrs)
}

func (p PointChangeRecordRepositoryImpl) Create(record *model.PointChangeRecord) {
  p.db.Create(record)
}

func (p PointChangeRecordRepositoryImpl) GetById(id uint64) (*model.PointChangeRecord, error) {
  var record model.PointChangeRecord
  err := p.db.First(&record, id).Error
  if err != nil {
    return nil, err
  }
  return &record, nil
}

func (p PointChangeRecordRepositoryImpl) Update(record *model.PointChangeRecord) {
  p.db.Updates(record)
}

func (p PointChangeRecordRepositoryImpl) Delete(id uint64) error {
  return p.db.Delete(&model.PointChangeRecord{}, id).Error
}

func (p PointChangeRecordRepositoryImpl) GetByNetworkId(networkId uint64) ([]model.PointChangeRecord, error) {
  var records []model.PointChangeRecord
  err := p.db.Where("network_id = ?", networkId).Find(&records).Error
  return records, err
}

func (p PointChangeRecordRepositoryImpl) GetByBlockNumber(blockNumber uint64) ([]model.PointChangeRecord, error) {
  var records []model.PointChangeRecord
  err := p.db.Where("block_number = ?", blockNumber).Find(&records).Error
  return records, err
}

func (p PointChangeRecordRepositoryImpl) GetByAddress(address string) ([]model.PointChangeRecord, error) {
  var records []model.PointChangeRecord
  err := p.db.Where("address = ?", address).Find(&records).Error
  return records, err
}

func (p PointChangeRecordRepositoryImpl) GetByNetworkAndBlock(networkId, blockNumber uint64) ([]model.PointChangeRecord, error) {
  var records []model.PointChangeRecord
  err := p.db.Where("network_id = ? AND block_number = ?", networkId, blockNumber).Find(&records).Error
  return records, err
}

func (p PointChangeRecordRepositoryImpl) GetByNetworkAndAddress(networkId uint64, address string) ([]model.PointChangeRecord, error) {
  var records []model.PointChangeRecord
  err := p.db.Where("network_id = ? AND address = ?", networkId, address).Find(&records).Error
  return records, err
}

func NewPointChangeRecordRepository(db *gorm.DB) PointChangeRecordRepository {
  return &PointChangeRecordRepositoryImpl{db: db}
}
