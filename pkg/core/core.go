package core

import (
	"errors"
	"fmt"
	"maps"
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
	CONFIG_KEY_NACOS    = "nacos"
)

var (
	defaultTimeZone = "Asia/Shanghai"
)

var defaultCore Core = NewContainer()

func Default() Core {
	return defaultCore
}

func SetDefault(c Core) {
	if c == nil {
		return
	}
	defaultCore = c
}

func Build(configFiles ...string) (*Container, error) {
	return BuildWithProvider(new(BaseAppConfig), configFiles...)
}

func BuildWithProvider(provider ConfigProvider, configFiles ...string) (*Container, error) {
	if len(configFiles) == 0 {
		return nil, errors.New("缺少配置文件")
	}
	if provider == nil {
		provider = new(BaseAppConfig)
	}
	if err := provider.LoadConfig(configFiles...); err != nil {
		return nil, fmt.Errorf("读取配置失败[%v]", err)
	}
	container, err := BuildWithConfig(provider.GetBaseAppConfig())
	if err != nil {
		return nil, err
	}
	if exporter, ok := provider.(CustomConfigExporter); ok {
		for key, cfg := range exporter.CustomConfigs() {
			container.SetCustomConfig(key, cfg)
		}
	}
	return container, nil
}

func BuildWithConfig(cfg *BaseAppConfig) (*Container, error) {
	container := NewContainer()
	if err := container.InitWithConfig(cfg); err != nil {
		return nil, err
	}
	SetDefault(container)
	return container, nil
}

type Container struct {
	appConfig *BaseAppConfig
	resources Resources

	customConfigs map[string]any

	mux sync.RWMutex
}

func NewContainer() *Container {
	return &Container{
		appConfig: &BaseAppConfig{
			System: &SysConfig{},
		},
		resources: Resources{
			loggers:   make(map[string]logger.Logger),
			databases: make(map[string]*gorm.DB),
			caches:    make(map[string]cache.AdapterCache),
			consumers: make(map[string][]rocketmq.Consumer),
		},
		customConfigs: make(map[string]any),
	}
}

func (e *Container) GetSystem() SysConfig {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return *e.appConfig.System
}

func (e *Container) SetServerId(serverId uint8) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.appConfig.System.ServerId = serverId
}

func (e *Container) GetServerId() uint8 {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.appConfig.System.ServerId
}

func (e *Container) SetDatabase(key string, db *gorm.DB) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.resources.setDatabase(key, db)
}

func (e *Container) GetDatabase(key string) *gorm.DB {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.resources.getDatabase(key)
}

func (e *Container) GetAllDatabases() map[string]*gorm.DB {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return maps.Clone(e.resources.databases)
}

func (e *Container) SetCustomConfig(key string, config any) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.customConfigs[key] = config
}

func (e *Container) GetCustomConfig(key string) any {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.customConfigs[key]
}

func (e *Container) SetLogger(key string, log logger.Logger) {
	e.mux.Lock()
	defer e.mux.Unlock()
	if key == DEFAULT_KEY {
		logger.SetLogger(log)
	}
	e.resources.setLogger(key, log)
}

func (e *Container) GetLogger(key string) logger.Logger {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.resources.getLogger(key)
}

func (e *Container) SyncLogger() {
	e.mux.Lock()
	defer e.mux.Unlock()
	for _, log := range e.resources.loggers {
		log.Sync()
	}
}

// SetCache 设置缓存
func (e *Container) SetCache(key string, c cache.AdapterCache) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.resources.setCache(key, c)
}

// Cache 获取缓存
func (e *Container) GetCache(key string) cache.AdapterCache {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.resources.getCache(key)
}

// SetLocker 设置分布式锁
func (e *Container) SetLocker(c locker.AdapterLocker) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.resources.setLocker(c)
}

// Locker 获取分布式锁
func (e *Container) GetLocker() locker.AdapterLocker {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.resources.getLocker()
}

func (e *Container) AddConsumer(key string, consumer rocketmq.Consumer) {
	e.mux.Lock()
	defer e.mux.Unlock()
	e.resources.addConsumer(key, consumer)
}

func (e *Container) GetConsumers(key string) []rocketmq.Consumer {
	e.mux.RLock()
	defer e.mux.RUnlock()
	return e.resources.getConsumer(key)
}

func (e *Container) InitWithConfig(cfg *BaseAppConfig) error {
	if cfg == nil {
		return errors.New("base app config is nil")
	}
	var err error
	if err = e.loadConfig(cfg); err != nil {
		return err
	}
	if err = e.setupLog(); err != nil {
		return err
	}
	if err = e.setupRedis(); err != nil {
		return err
	}
	if err = e.setupDB(); err != nil {
		return err
	}
	if err = e.setServerId(); err != nil {
		return err
	}
	return nil
}

type Resources struct {
	loggers   map[string]logger.Logger
	databases map[string]*gorm.DB
	caches    map[string]cache.AdapterCache
	consumers map[string][]rocketmq.Consumer
	locker    locker.AdapterLocker
}

func (r *Resources) setLogger(key string, log logger.Logger) {
	r.loggers[key] = log
}

func (r *Resources) getLogger(key string) logger.Logger {
	return r.loggers[key]
}

func (r *Resources) setDatabase(key string, db *gorm.DB) {
	r.databases[key] = db
}

func (r *Resources) getDatabase(key string) *gorm.DB {
	return r.databases[key]
}

func (r *Resources) setCache(key string, cache cache.AdapterCache) {
	r.caches[key] = cache
}

func (r *Resources) getCache(key string) cache.AdapterCache {
	return r.caches[key]
}

func (r *Resources) addConsumer(key string, consumer rocketmq.Consumer) {
	r.consumers[key] = append(r.consumers[key], consumer)
}

func (r *Resources) getConsumer(key string) []rocketmq.Consumer {
	return append([]rocketmq.Consumer{}, r.consumers[key]...)
}

func (r *Resources) setLocker(locker locker.AdapterLocker) {
	r.locker = locker
}

func (r *Resources) getLocker() locker.AdapterLocker {
	return r.locker
}
