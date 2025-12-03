package repository

import (
  "com.wh1200.points/internal/model"
  "gorm.io/gorm"
)

type TransferRecordRepository interface {
  Create(record *model.TransferRecord)
  GetById(id uint64) (*model.TransferRecord, error)
  Update(record *model.TransferRecord)
  Delete(id uint64) error
  GetByHash(hash string) (*model.TransferRecord, error)
  GetByNetworkId(networkId uint64) ([]model.TransferRecord, error)
  GetByBlockNumber(blockNumber uint64) ([]model.TransferRecord, error)
  GetByFromAddress(fromAddress string) ([]model.TransferRecord, error)
  GetByToAddress(toAddress string) ([]model.TransferRecord, error)
  GetByNetworkAndBlock(networkId, blockNumber uint64) ([]model.TransferRecord, error)
  SaveAll(records []model.TransferRecord)
  FindExistsHashes(hashes []string) []string
  GetBalanceOfAll(networkId uint64) []Balance
  GetAllNotCalc(networkId uint64) []model.TransferRecord
}

type TransferRecordRepositoryImpl struct {
  db *gorm.DB
}

type Balance struct {
  Address string
  Balance uint64
}

func (t TransferRecordRepositoryImpl) GetAllNotCalc(networkId uint64) []model.TransferRecord {
  records := make([]model.TransferRecord, 0)
  t.db.Model(&model.TransferRecord{}).Where("network_id = ? and calc_status = ?", networkId, 0).Find(&records)
  return records
}
func (t TransferRecordRepositoryImpl) GetBalanceOfAll(networkId uint64) []Balance {
  incomingBalance := make([]Balance, 0)
  outgoingBalance := make([]Balance, 0)
  t.db.Model(&model.TransferRecord{}).Select("sum(value) as balance", "from_address as address").
    Where("network_id = ?", networkId).
    Group("from_address").Find(&outgoingBalance)
  t.db.Model(&model.TransferRecord{}).Select("sum(value) as balance", "to_address as address").
    Where("network_id = ?", networkId).
    Group("to_address").Find(&incomingBalance)
  outgoingMap := make(map[string]uint64)
  for _, balance := range outgoingBalance {
    outgoingMap[balance.Address] = balance.Balance
  }
  res := make([]Balance, 0)
  for _, incoming := range incomingBalance {
    in := incoming.Balance
    out := outgoingMap[incoming.Address]
    res = append(res, Balance{Address: incoming.Address, Balance: in - out})
  }
  return res
}

func (t TransferRecordRepositoryImpl) SaveAll(records []model.TransferRecord) {
  if len(records) == 0 {
    return
  }
  hashes := make([]string, len(records))
  for i, record := range records {
    hashes[i] = record.Hash
  }
  existsHashes := t.FindExistsHashes(hashes)
  // 3. 用 map 做快速过滤
  existsMap := make(map[string]struct{}, len(existsHashes))
  for _, h := range existsHashes {
    existsMap[h] = struct{}{}
  }

  // 4. 过滤掉已存在的记录
  newRecords := make([]model.TransferRecord, 0, len(records))
  for _, record := range records {
    if _, exists := existsMap[record.Hash]; !exists {
      newRecords = append(newRecords, record)
    }
  }

  // 5. 批量插入
  if len(newRecords) > 0 {
    err := t.db.Create(&newRecords).Error
    if err != nil {
      panic(err)
    }
  }
}

func (t TransferRecordRepositoryImpl) FindExistsHashes(hashes []string) []string {
  var res []string
  t.db.Select("hash").Where("hash IN (?)", hashes).Find(&res)
  return res
}

func (t TransferRecordRepositoryImpl) Create(record *model.TransferRecord) {
  t.db.Create(record)
}

func (t TransferRecordRepositoryImpl) GetById(id uint64) (*model.TransferRecord, error) {
  var record model.TransferRecord
  err := t.db.First(&record, id).Error
  if err != nil {
    return nil, err
  }
  return &record, nil
}

func (t TransferRecordRepositoryImpl) Update(record *model.TransferRecord) {
  t.db.Updates(record)
}

func (t TransferRecordRepositoryImpl) Delete(id uint64) error {
  return t.db.Delete(&model.TransferRecord{}, id).Error
}

func (t TransferRecordRepositoryImpl) GetByHash(hash string) (*model.TransferRecord, error) {
  var record model.TransferRecord
  err := t.db.Where("hash = ?", hash).First(&record).Error
  if err != nil {
    return nil, err
  }
  return &record, nil
}

func (t TransferRecordRepositoryImpl) GetByNetworkId(networkId uint64) ([]model.TransferRecord, error) {
  var records []model.TransferRecord
  err := t.db.Where("network_id = ?", networkId).Find(&records).Error
  return records, err
}

func (t TransferRecordRepositoryImpl) GetByBlockNumber(blockNumber uint64) ([]model.TransferRecord, error) {
  var records []model.TransferRecord
  err := t.db.Where("block_number = ?", blockNumber).Find(&records).Error
  return records, err
}

func (t TransferRecordRepositoryImpl) GetByFromAddress(fromAddress string) ([]model.TransferRecord, error) {
  var records []model.TransferRecord
  err := t.db.Where("from_address = ?", fromAddress).Find(&records).Error
  return records, err
}

func (t TransferRecordRepositoryImpl) GetByToAddress(toAddress string) ([]model.TransferRecord, error) {
  var records []model.TransferRecord
  err := t.db.Where("to_address = ?", toAddress).Find(&records).Error
  return records, err
}

func (t TransferRecordRepositoryImpl) GetByNetworkAndBlock(networkId, blockNumber uint64) ([]model.TransferRecord, error) {
  var records []model.TransferRecord
  err := t.db.Where("network_id = ? AND block_number = ?", networkId, blockNumber).Find(&records).Error
  return records, err
}

func NewTransferRecordRepository(db *gorm.DB) TransferRecordRepository {
  return &TransferRecordRepositoryImpl{db: db}
}
