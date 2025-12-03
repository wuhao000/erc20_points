package repository

import (
  "com.wh1200.points/internal/model"
  "gorm.io/gorm"
)

type NetworkRepository interface {
  Create(network *model.Network)
  GetById(id uint64) (*model.Network, error)
  Update(network *model.Network)
  Delete(id uint64) error
  GetByName(name string) (*model.Network, error)
  GetAll() ([]model.Network, error)
}

type NetworkRepositoryImpl struct {
  db *gorm.DB
}

func (n NetworkRepositoryImpl) Create(network *model.Network) {
  n.db.Create(network)
}

func (n NetworkRepositoryImpl) GetById(id uint64) (*model.Network, error) {
  var network model.Network
  err := n.db.First(&network, id).Error
  if err != nil {
    return nil, err
  }
  return &network, nil
}

func (n NetworkRepositoryImpl) Update(network *model.Network) {
  n.db.Updates(network)
}

func (n NetworkRepositoryImpl) Delete(id uint64) error {
  return n.db.Delete(&model.Network{}, id).Error
}

func (n NetworkRepositoryImpl) GetByName(name string) (*model.Network, error) {
  var network model.Network
  err := n.db.Where("name = ?", name).First(&network).Error
  if err != nil {
    return nil, err
  }
  return &network, nil
}

func (n NetworkRepositoryImpl) GetAll() ([]model.Network, error) {
  var networks []model.Network
  err := n.db.Find(&networks).Error
  return networks, err
}

func NewNetworkRepository(db *gorm.DB) NetworkRepository {
  return &NetworkRepositoryImpl{db: db}
}