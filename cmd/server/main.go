package main

import (
  "flag"

  "com.wh1200.points/internal/config"
  "com.wh1200.points/internal/service"
)

const (
  defaultConfigPath = "./config/config.yml"
)

func main() {
  conf := flag.String("config", defaultConfigPath, "config file path")
  flag.Parse()
  err := config.LoadConfig(*conf)
  if err != nil {
    panic(err)
  }
  db := config.InitDb()
  service.InjectRepositories(db)
  service.Start()
  select {} // 阻塞主线程，程序不会退出
}
