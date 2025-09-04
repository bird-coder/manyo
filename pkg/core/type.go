/*
 * @Author: yujiajie
 * @Date: 2024-12-25 19:42:41
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-09-03 18:13:25
 * @FilePath: /manyo/pkg/core/type.go
 * @Description:
 */
package core

import (
	"github.com/bird-coder/manyo/lib/rocketmq"
	"github.com/bird-coder/manyo/pkg/logger"
	"github.com/bird-coder/manyo/pkg/storage/cache"
	"github.com/bird-coder/manyo/pkg/storage/locker"

	"gorm.io/gorm"
)

type Core interface {
	SetSysInfo(cfg *SysConfig)
	GetSysInfo() *SysConfig

	SetServerId(serverId uint8)
	GetServerId() uint8

	SetDb(key string, db *gorm.DB)
	GetDb(key string) *gorm.DB
	GetAllDb() map[string]*gorm.DB

	SetLogger(key string, log logger.Logger)
	GetLogger(key string) logger.Logger
	SyncLogger()

	SetConfig(key string, config any)
	GetConfig(key string) any

	SetCacheAdapter(key string, c cache.AdapterCache)
	GetCacheAdapter(key string) cache.AdapterCache

	SetLockerAdapter(locker.AdapterLocker)
	GetLockerAdapter() locker.AdapterLocker

	AddConsumer(key string, consumer rocketmq.Consumer)
	GetConsumer(key string) []rocketmq.Consumer

	Init(string) error
}
