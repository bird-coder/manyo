/*
 * @Author: yujiajie
 * @Date: 2024-12-25 20:13:05
 * @LastEditors: yujiajie
 * @LastEditTime: 2024-12-27 16:24:54
 * @FilePath: /manyo/pkg/core/initialize.go
 * @Description:
 */
package core

import (
	"context"
	"fmt"

	"github.com/bird-coder/manyo/config"
	"github.com/bird-coder/manyo/pkg/logger"
	"github.com/bird-coder/manyo/pkg/storage/cache"
	"github.com/bird-coder/manyo/pkg/storage/locker"
)

const (
	SERVER_ID = "nbgame:rock:server_id:%s"
)

func setupRedis() error {
	rdsConfigs := Kernal.GetConfig(CONFIG_KEY_REDIS).(map[string]*config.RedisDailConfig)
	for k, cfg := range rdsConfigs {
		rds, err := cache.NewRedis(nil, cfg)
		if err != nil {
			return fmt.Errorf("cache setup error, %v", err)
		}
		Kernal.SetCacheAdapter(k, rds)
	}
	lockConfig := Kernal.GetConfig(CONFIG_KEY_LOCKER).(*config.RedisDailConfig)
	if lockConfig != nil {
		r, err := locker.NewRedis(nil, lockConfig)
		if err != nil {
			return fmt.Errorf("locker setup error, %v", err)
		}
		Kernal.SetLockerAdapter(r)
	}
	return nil
}

func setupDB() error {
	return nil
}

func setupLog() {
	logConfigs := Kernal.GetConfig(CONFIG_KEY_LOGGER).(map[string]*config.LoggerConfig)
	for k, cfg := range logConfigs {
		log := logger.NewLogger(cfg, Kernal.GetEnv())
		Kernal.SetLogger(k, log)
	}
}

func setServerId() error {
	rds := Kernal.GetCacheAdapter(DEFAULT_KEY).(*cache.Redis)
	res, err := rds.GetClient().Incr(context.TODO(), fmt.Sprintf(SERVER_ID, Kernal.GetAppName())).Result()
	if err != nil {
		return fmt.Errorf("生成serverid失败, err: %v", err)
	}
	Kernal.SetServerId(uint8(res))
	return nil
}
