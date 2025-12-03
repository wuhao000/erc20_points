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
  // 初始化数据库，使用postgresql
  db := config.InitDb()
  service.InjectRepositories(db)
  // 启动事件同步以及积分计算任务
  service.Start()
  select {} // 阻塞主线程，程序不会退出
}
