package config

import (
  "github.com/spf13/viper"
  "github.com/zeromicro/go-zero/core/stores/redis"
)

type ChainConfig struct {
  Name            string
  DelayBlocks     uint64
  StartBlock      uint64
  ChainId         uint64
  NodeUrl         string
  NodeUrl2        string
  PollingInterval uint64 // 轮询间隔，单位：秒
  ContractAddress string
}

type Config struct {
  NodeApiKey string
  Chains     []ChainConfig
  Events     []string
  Redis      *redis.RedisConf
  Database   Database
}

type Database struct {
  Dialect  string
  Database string
  Username string
  Password string
  Host     string
  Port     int
}

var Cfg *Config

func LoadConfig(configFilePath string) error {
  viper.SetConfigFile(configFilePath)
  viper.SetConfigType("yaml")
  viper.AutomaticEnv()
  if err := viper.ReadInConfig(); err != nil {
    return err
  }
  config := &Config{}
  if err := viper.Unmarshal(config); err != nil {
    return err
  }
  Cfg = config
  // Implementation for loading configuration
  return nil
}
