/*
 * @Author: yujiajie
 * @Date: 2024-12-25 19:42:41
 * @LastEditors: yujiajie
 * @LastEditTime: 2024-12-27 16:25:09
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
	SetAppName(name string)
	GetAppName() string

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

	SetEnv(env string)
	GetEnv() string

	SetServerId(serverId uint8)
	GetServerId() uint8

	Init() error
}
