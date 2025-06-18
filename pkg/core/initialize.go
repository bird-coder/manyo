/*
 * @Author: yujiajie
 * @Date: 2024-12-25 20:13:05
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-06-04 18:11:37
 * @FilePath: /manyo/pkg/core/initialize.go
 * @Description:
 */
package core

import (
	"fmt"
	"math"

	"github.com/bird-coder/manyo/config"
	"github.com/bird-coder/manyo/pkg/generator"
	"github.com/bird-coder/manyo/pkg/logger"
	"github.com/bird-coder/manyo/pkg/storage/cache"
	"github.com/bird-coder/manyo/pkg/storage/database"
	"github.com/bird-coder/manyo/pkg/storage/locker"
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
	dbConfigs := Kernal.GetConfig(CONFIG_KEY_DATABASE).(map[string]*config.MysqlConfig)
	for k, cfg := range dbConfigs {
		db, err := database.NewMysql(k, cfg)
		if err != nil {
			return fmt.Errorf("db setup error, %v", err)
		}
		Kernal.SetDb(k, db)
	}
	return nil
}

func setupLog() {
	logConfigs := Kernal.GetConfig(CONFIG_KEY_LOGGER).(map[string]*config.LoggerConfig)
	for k, cfg := range logConfigs {
		log := logger.NewLogger(cfg, Kernal.GetSysInfo().Environment)
		Kernal.SetLogger(k, log)
	}
}

func setServerId() error {
	//如果serverId通过手动配置，则使用配置，否则通过redis自增id生成
	if Kernal.GetServerId() > 0 {
		return nil
	}
	rds := Kernal.GetCacheAdapter(DEFAULT_KEY)
	if rds == nil {
		return fmt.Errorf("生成serverid失败, cache未初始化")
	}
	redisCli := rds.(*cache.Redis)
	generator := generator.NewRedisGenerator(redisCli.GetClient(), math.MaxUint8, generator.WithAppName(Kernal.GetSysInfo().AppName))
	res := generator.Next()
	if res == 0 {
		return fmt.Errorf("生成serverid失败, res: %d", res)
	}
	Kernal.SetServerId(uint8(res))
	return nil
}
