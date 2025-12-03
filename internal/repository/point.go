package repository

import (
  "com.wh1200.points/internal/model"
  "gorm.io/gorm"
  "gorm.io/gorm/clause"
)

type PointRecordRepository interface {
  Create(point *model.PointRecord)
  GetById(id uint64) (*model.PointRecord, error)
  Update(point *model.PointRecord)
  Delete(id uint64) error
  GetByUserId(userId uint64) ([]model.PointRecord, error)
  GetByNetworkId(networkId uint64) ([]model.PointRecord, error)
  GetByAddress(address string) (*model.PointRecord, error)
  FindByNetworkId(id uint64) []model.PointRecord
  SaveOrUpdateAll(prs []*model.PointRecord)
  FindAddressesByNetworkId(networkId uint64) []string
}

type PointRecordRepositoryImpl struct {
  db *gorm.DB
}

func (p PointRecordRepositoryImpl) FindAddressesByNetworkId(networkId uint64) []string {
  var addresses []string
  p.db.Model(&model.PointRecord{}).Select("address").Where("network_id = ?", networkId).Find(&addresses)
  return addresses
}

func (p PointRecordRepositoryImpl) SaveOrUpdateAll(prs []*model.PointRecord) {
  if len(prs) == 0 {
    return
  }
  p.db.Clauses(clause.OnConflict{
    Columns: []clause.Column{{Name: "id"}}, // 冲突键
    DoUpdates: clause.AssignmentColumns([]string{
      "balance", "points", "last_updated",
    }),
  }).Create(&prs)
}

func (p PointRecordRepositoryImpl) FindByNetworkId(id uint64) []model.PointRecord {
  records := make([]model.PointRecord, 0)
  p.db.Model(&model.PointRecord{}).Where("network_id = ?", id).Find(&records)
  return records
}

func (p PointRecordRepositoryImpl) Create(point *model.PointRecord) {
  p.db.Create(point)
}

func (p PointRecordRepositoryImpl) GetById(id uint64) (*model.PointRecord, error) {
  var point model.PointRecord
  err := p.db.First(&point, id).Error
  if err != nil {
    return nil, err
  }
  return &point, nil
}

func (p PointRecordRepositoryImpl) Update(point *model.PointRecord) {
  p.db.Updates(point)
}

func (p PointRecordRepositoryImpl) Delete(id uint64) error {
  return p.db.Delete(&model.PointRecord{}, id).Error
}

func (p PointRecordRepositoryImpl) GetByUserId(userId uint64) ([]model.PointRecord, error) {
  var points []model.PointRecord
  err := p.db.Where("user_id = ?", userId).Find(&points).Error
  return points, err
}

func (p PointRecordRepositoryImpl) GetByNetworkId(networkId uint64) ([]model.PointRecord, error) {
  var points []model.PointRecord
  err := p.db.Where("network_id = ?", networkId).Find(&points).Error
  return points, err
}

func (p PointRecordRepositoryImpl) GetByAddress(address string) (*model.PointRecord, error) {
  var point model.PointRecord
  err := p.db.Where("address = ?", address).First(&point).Error
  if err != nil {
    return nil, err
  }
  return &point, nil
}

func NewPointRecordRepository(db *gorm.DB) PointRecordRepository {
  return &PointRecordRepositoryImpl{db: db}
}
