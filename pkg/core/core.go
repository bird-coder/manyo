package core

import (
	"sync"

	"github.com/bird-coder/manyo/lib/rocketmq"
	"github.com/bird-coder/manyo/pkg/logger"
	"github.com/bird-coder/manyo/pkg/storage/cache"
	"github.com/bird-coder/manyo/pkg/storage/locker"

	"gorm.io/gorm"
)

const (
	DEFAULT_KEY = "default"

	CONFIG_KEY_SERVER   = "server"
	CONFIG_KEY_TARGET   = "targets"
	CONFIG_KEY_LOGGER   = "loggers"
	CONFIG_KEY_REDIS    = "redis"
	CONFIG_KEY_DATABASE = "databases"
	CONFIG_KEY_LOCKER   = "locker"
	CONFIG_KEY_ROCKET   = "rocketmq"
	CONFIG_KEY_CONSUMER = "consumers"
	CONFIG_KEY_DISCOVER = "discover"
)

var Kernal Core = NewKernal()

type Container struct {
	name      string
	dbs       map[string]*gorm.DB
	configs   map[string]any
	logs      map[string]logger.Logger
	rds       map[string]cache.AdapterCache
	consumers map[string][]rocketmq.Consumer
	locker    locker.AdapterLocker
	env       string
	serverId  uint8

	mux sync.RWMutex
}

func NewKernal() *Container {
	return &Container{
		dbs:       make(map[string]*gorm.DB),
		configs:   make(map[string]any),
		logs:      make(map[string]logger.Logger),
		rds:       make(map[string]cache.AdapterCache),
		consumers: make(map[string][]rocketmq.Consumer),
	}
}

func (e *Container) SetAppName(name string) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.name = name
}

func (e *Container) GetAppName() string {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.name
}

func (e *Container) SetDb(key string, db *gorm.DB) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.dbs[key] = db
}

func (e *Container) GetDb(key string) *gorm.DB {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.dbs[key]
}

func (e *Container) GetAllDb() map[string]*gorm.DB {
	return e.dbs
}

func (e *Container) SetConfig(key string, config any) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.configs[key] = config
}

func (e *Container) GetConfig(key string) any {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.configs[key]
}

func (e *Container) SetLogger(key string, log logger.Logger) {
	e.mux.Lock()
	defer e.mux.Unlock()
	if key == DEFAULT_KEY {
		logger.SetLogger(log)
	}
	e.logs[key] = log
}

func (e *Container) GetLogger(key string) logger.Logger {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.logs[key]
}

func (e *Container) SyncLogger() {
	e.mux.Lock()
	defer e.mux.Unlock()
	for _, log := range e.logs {
		log.Sync()
	}
}

// SetCacheAdapter 设置缓存
func (e *Container) SetCacheAdapter(key string, c cache.AdapterCache) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.rds[key] = c
}

// GetCacheAdapter 获取缓存
func (e *Container) GetCacheAdapter(key string) cache.AdapterCache {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.rds[key]
}

// SetLockerAdapter 设置分布式锁
func (e *Container) SetLockerAdapter(c locker.AdapterLocker) {
	e.locker = c
}

// GetLockerAdapter 获取分布式锁
func (e *Container) GetLockerAdapter() locker.AdapterLocker {
	return e.locker
}

func (e *Container) AddConsumer(key string, consumer rocketmq.Consumer) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.consumers[key] = append(e.consumers[key], consumer)
}

func (e *Container) GetConsumer(key string) []rocketmq.Consumer {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.consumers[key]
}

func (e *Container) SetEnv(env string) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.env = env
}

func (e *Container) GetEnv() string {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.env
}

func (e *Container) SetServerId(serverId uint8) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.serverId = serverId
}

func (e *Container) GetServerId() uint8 {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.serverId
}

func (e *Container) Init() error {
	var err error
	setupLog()
	if err = setupRedis(); err != nil {
		return err
	}
	if err = setupDB(); err != nil {
		return err
	}
	if err = setServerId(); err != nil {
		return err
	}
	return nil
}
