package service

import "com.wh1200.points/internal/config"

func Start() {
  for _, chainConfig := range config.Cfg.Chains {
    // 启动事件同步任务
    startSyncTask(&chainConfig)
    // 启动积分计算任务
    startCalcPointsTask(&chainConfig)
  }
}
