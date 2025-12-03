package repository

import (
  "com.wh1200.points/internal/model"
  "gorm.io/gorm"
)

type UserRepository interface {
  Create(user *model.User)
  GetById(id uint64) (*model.User, error)
  Update(user *model.User)
}

type UserRepositoryImpl struct {
  db *gorm.DB
}

func (u UserRepositoryImpl) Create(user *model.User) {
  u.db.Create(user)
}

func (u UserRepositoryImpl) GetById(id uint64) (*model.User, error) {
  var user model.User
  err := u.db.First(&user, id).Error
  if err != nil {
    panic(err)
  }
  return &user, nil
}

func (u UserRepositoryImpl) Update(user *model.User) {
  u.db.Updates(user)
}

func NewUserRepository(db *gorm.DB) UserRepository {
  return &UserRepositoryImpl{db: db}
}
