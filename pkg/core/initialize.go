/*
 * @Author: yujiajie
 * @Date: 2024-12-25 20:13:05
 * @LastEditors: yujiajie 1037297660@qq.com
 * @LastEditTime: 2026-05-11 18:40:47
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

func loadConfig(configFile string) error {
	appConfig := new(BaseAppConfig)
	if err := appConfig.LoadConfig(configFile); err != nil {
		return fmt.Errorf("读取配置失败[%v]", err)
	}
	return nil
}

func getConfig[T any](key string) (T, error) {
	var cfg T
	v := Kernal.GetConfig(key)
	if v == nil {
		return cfg, fmt.Errorf("config %s is nil", key)
	}
	res, ok := v.(T)
	if !ok {
		return cfg, fmt.Errorf("config %s has wrong type", key)
	}
	return res, nil
}

func setupRedis() error {
	rdsConfigs, err := getConfig[map[string]*config.RedisDailConfig](CONFIG_KEY_REDIS)
	if err != nil {
		return fmt.Errorf("cache setup error, %v", err)
	}
	for k, cfg := range rdsConfigs {
		rds, err := cache.NewRedis(nil, cfg)
		if err != nil {
			return fmt.Errorf("cache setup error, %v", err)
		}
		Kernal.SetCacheAdapter(k, rds)
	}
	lockConfig, err := getConfig[*config.RedisDailConfig](CONFIG_KEY_LOCKER)
	if err != nil {
		return fmt.Errorf("locker setup error, %v", err)
	}
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
	dbConfigs, err := getConfig[map[string]*config.MysqlConfig](CONFIG_KEY_DATABASE)
	if err != nil {
		return fmt.Errorf("db setup error, %v", err)
	}
	for k, cfg := range dbConfigs {
		db, err := database.NewMysql(k, cfg)
		if err != nil {
			return fmt.Errorf("db setup error, %v", err)
		}
		Kernal.SetDb(k, db)
	}
	return nil
}

func setupLog() error {
	logConfigs, err := getConfig[map[string]*config.LoggerConfig](CONFIG_KEY_LOGGER)
	if err != nil {
		return fmt.Errorf("log setup error, %v", err)
	}
	for k, cfg := range logConfigs {
		log := logger.NewLogger(cfg, Kernal.GetSysInfo().Environment)
		Kernal.SetLogger(k, log)
	}
	return nil
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
