/*
 * @Author: yujiajie
 * @Date: 2024-12-25 19:42:41
 * @LastEditors: yujiajie 1037297660@qq.com
 * @LastEditTime: 2026-05-12 11:45:38
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
	GetSystem() SysConfig

	GetServerId() uint8

	GetDatabase(key string) *gorm.DB
	GetAllDatabases() map[string]*gorm.DB

	GetLogger(key string) logger.Logger
	SyncLogger()

	SetCustomConfig(key string, config any)
	GetCustomConfig(key string) any

	GetCache(key string) cache.AdapterCache

	GetLocker() locker.AdapterLocker

	GetConsumers(key string) []rocketmq.Consumer
}
