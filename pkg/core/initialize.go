/*
 * @Author: yujiajie
 * @Date: 2024-12-25 20:13:05
 * @LastEditors: yujiajie 1037297660@qq.com
 * @LastEditTime: 2026-05-14 09:59:11
 * @FilePath: /manyo/pkg/core/initialize.go
 * @Description:
 */
package core

import (
	"fmt"
	"math"

	"github.com/bird-coder/manyo/pkg/generator"
	"github.com/bird-coder/manyo/pkg/logger"
	"github.com/bird-coder/manyo/pkg/storage/cache"
	"github.com/bird-coder/manyo/pkg/storage/database"
	"github.com/bird-coder/manyo/pkg/storage/locker"
)

type ConfigReader interface {
	GetCustomConfig(key string) any
}

func GetCustomConfig[T any](r ConfigReader, key string) (T, error) {
	var cfg T
	v := r.GetCustomConfig(key)
	if v == nil {
		return cfg, fmt.Errorf("config %s is nil", key)
	}
	res, ok := v.(T)
	if !ok {
		return cfg, fmt.Errorf("config %s has wrong type", key)
	}
	return res, nil
}

func (c *Container) loadConfig(cfg *BaseAppConfig) error {
	cfg.Normalize()
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("配置未通过校验[%v]", err)
	}
	c.appConfig = cfg
	return nil
}

func (c *Container) setupRedis() error {
	rdsConfigs := c.appConfig.Redis
	for k, cfg := range rdsConfigs {
		rds, err := cache.NewRedis(nil, cfg)
		if err != nil {
			return fmt.Errorf("cache setup error, %v", err)
		}
		c.SetCache(k, rds)
	}
	lockConfig := c.appConfig.Locker
	if lockConfig != nil {
		r, err := locker.NewRedis(nil, lockConfig)
		if err != nil {
			return fmt.Errorf("locker setup error, %v", err)
		}
		c.SetLocker(r)
	}
	return nil
}

func (c *Container) setupDB() error {
	dbConfigs := c.appConfig.Databases
	for k, cfg := range dbConfigs {
		db, err := database.NewMysql(k, cfg)
		if err != nil {
			return fmt.Errorf("db setup error, %v", err)
		}
		c.SetDatabase(k, db)
	}
	return nil
}

func (c *Container) setupLog() error {
	logConfigs := c.appConfig.Loggers
	for k, cfg := range logConfigs {
		log := logger.NewLogger(cfg, c.GetSystem().Environment)
		c.SetLogger(k, log)
	}
	return nil
}

func (c *Container) setServerId() error {
	//如果serverId通过手动配置，则使用配置，否则通过redis自增id生成
	if c.GetServerId() > 0 {
		return nil
	}
	rds := c.GetCache(DEFAULT_KEY)
	if rds == nil {
		return fmt.Errorf("生成serverid失败, cache未初始化")
	}
	redisCli := rds.(*cache.Redis)
	generator := generator.NewRedisGenerator(redisCli.GetClient(), math.MaxUint8, generator.WithAppName(c.GetSystem().AppName))
	res := generator.Next()
	if res == 0 {
		return fmt.Errorf("生成serverid失败, res: %d", res)
	}
	c.SetServerId(uint8(res))
	return nil
}
