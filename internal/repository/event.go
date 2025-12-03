package repository

import (
  "com.wh1200.points/internal/model"
  "gorm.io/gorm"
)

type EventRepository interface {
  Create(event *model.Event)
  GetById(id string) (*model.Event, error)
  Update(event *model.Event)
  Delete(id string) error
  GetByChainId(chainId uint64) ([]model.Event, error)
  GetByNetworkId(networkId uint64) ([]model.Event, error)
  GetByBlockNumber(blockNumber uint64) ([]model.Event, error)
  GetByNetworkAndBlock(networkId, blockNumber uint64) ([]model.Event, error)
  GetByName(name string) ([]model.Event, error)
  SaveAll(events []model.Event)
  FindExistsHashes(hashes []string) ([]string, error)
}

type EventRepositoryImpl struct {
  db *gorm.DB
}

func (e EventRepositoryImpl) SaveAll(events []model.Event) {
  if len(events) == 0 {
    return
  }
  hashes := make([]string, len(events))
  for i, event := range events {
    hashes[i] = event.Hash
  }
  existsHashes, err := e.FindExistsHashes(hashes)
  if err != nil {
    panic(err)
  }
  // 3. 用 map 做快速过滤
  existsMap := make(map[string]struct{}, len(existsHashes))
  for _, h := range existsHashes {
    existsMap[h] = struct{}{}
  }

  // 4. 过滤掉已存在的事件
  newEvents := make([]model.Event, 0, len(events))
  for _, event := range events {
    if _, exists := existsMap[event.Hash]; !exists {
      newEvents = append(newEvents, event)
    }
  }

  // 5. 批量插入
  if len(newEvents) > 0 {
    err := e.db.Create(&newEvents).Error
    if err != nil {
      panic(err)
    }
  }

}

func (e EventRepositoryImpl) FindExistsHashes(hashes []string) ([]string, error) {
  var res []string
  err := e.db.Model(&model.Event{}).Select("hash").Where("hash IN (?)", hashes).Find(&res).Error
  return res, err
}

func (e EventRepositoryImpl) Create(event *model.Event) {
  e.db.Create(event)
}

func (e EventRepositoryImpl) GetById(id string) (*model.Event, error) {
  var event model.Event
  err := e.db.First(&event, "id = ?", id).Error
  if err != nil {
    return nil, err
  }
  return &event, nil
}

func (e EventRepositoryImpl) Update(event *model.Event) {
  e.db.Updates(event)
}

func (e EventRepositoryImpl) Delete(id string) error {
  return e.db.Delete(&model.Event{}, "id = ?", id).Error
}

func (e EventRepositoryImpl) GetByChainId(chainId uint64) ([]model.Event, error) {
  var events []model.Event
  err := e.db.Where("chain_id = ?", chainId).Find(&events).Error
  return events, err
}

func (e EventRepositoryImpl) GetByNetworkId(networkId uint64) ([]model.Event, error) {
  var events []model.Event
  err := e.db.Where("network_id = ?", networkId).Find(&events).Error
  return events, err
}

func (e EventRepositoryImpl) GetByBlockNumber(blockNumber uint64) ([]model.Event, error) {
  var events []model.Event
  err := e.db.Where("block_number = ?", blockNumber).Find(&events).Error
  return events, err
}

func (e EventRepositoryImpl) GetByNetworkAndBlock(networkId, blockNumber uint64) ([]model.Event, error) {
  var events []model.Event
  err := e.db.Where("network_id = ? AND block_number = ?", networkId, blockNumber).Find(&events).Error
  return events, err
}

func (e EventRepositoryImpl) GetByName(name string) ([]model.Event, error) {
  var events []model.Event
  err := e.db.Where("name = ?", name).Find(&events).Error
  return events, err
}

func NewEventRepository(db *gorm.DB) EventRepository {
  return &EventRepositoryImpl{db: db}
}
