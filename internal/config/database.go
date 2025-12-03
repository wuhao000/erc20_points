package config

import (
  "fmt"

  "com.wh1200.points/internal/model"
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
)

func InitDb() *gorm.DB {
  dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
    Cfg.Database.Host,
    Cfg.Database.Port,
    Cfg.Database.Username,
    Cfg.Database.Password,
    Cfg.Database.Database,
  )
  db, err := gorm.Open(
    postgres.Open(dsn),
    &gorm.Config{
      PrepareStmt: true,
    },
  )
  if err != nil {
    panic(err)
  }

  models := []interface{}{
    &model.User{},
    &model.Event{},
    &model.Network{},
    &model.PointRecord{},
    &model.PointChangeRecord{},
    &model.TransferRecord{},
  }

  for _, m := range models {
    err := db.AutoMigrate(m)
    if err != nil {
      panic(err)
    }
  }
  return db
}
