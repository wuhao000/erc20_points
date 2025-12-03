package service

import (
  "log"

  "com.wh1200.points/internal/config"
  "github.com/zeromicro/go-zero/core/stores/cache"
  "github.com/zeromicro/go-zero/core/stores/kv"
)

var store kv.Store

func getKvStore() kv.Store {
  if store != nil {
    return store
  }
  log.Printf("redis config: %+v\n", *config.Cfg.Redis)
  c := []cache.NodeConf{
    {
      RedisConf: *config.Cfg.Redis,
      Weight:    100,
    },
  }
  if len(c) == 0 || cache.TotalWeights(c) <= 0 {
    log.Fatal("no cache nodes")
  }
  store = kv.NewStore(c)
  return store
}
