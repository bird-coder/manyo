package core

import (
	"fmt"
	"maps"
	"slices"

	"github.com/bird-coder/manyo/lib/rocketmq"
	"github.com/bird-coder/manyo/pkg/logger"
	"github.com/bird-coder/manyo/pkg/storage/cache"
	"gorm.io/gorm"
)

// LoggerRegistry 提供对日志实例集合的只读访问。
type LoggerRegistry struct {
	logs map[string]logger.Logger
}

func NewLoggerRegistry(logs map[string]logger.Logger) *LoggerRegistry {
	return &LoggerRegistry{
		logs: logs,
	}
}

func (r *LoggerRegistry) Get(key string) (logger.Logger, bool) {
	if r == nil {
		return nil, false
	}
	log, ok := r.logs[key]
	return log, ok
}

func (r *LoggerRegistry) Default() (logger.Logger, bool) {
	return r.Get(DEFAULT_KEY)
}

func (r *LoggerRegistry) MustGet(key string) logger.Logger {
	log, ok := r.Get(key)
	if !ok {
		panic(fmt.Sprintf("logger %q not found", key))
	}
	return log
}

func (r *LoggerRegistry) All() map[string]logger.Logger {
	if r == nil {
		return nil
	}
	return maps.Clone(r.logs)
}

func (r *LoggerRegistry) Keys() []string {
	if r == nil {
		return nil
	}
	keys := maps.Keys(r.logs)
	return slices.Collect(keys)
}

// DatabaseRegistry 提供对数据库实例集合的只读访问。
type DatabaseRegistry struct {
	dbs map[string]*gorm.DB
}

func NewDatabaseRegistry(dbs map[string]*gorm.DB) *DatabaseRegistry {
	return &DatabaseRegistry{
		dbs: dbs,
	}
}

func (r *DatabaseRegistry) Get(key string) (*gorm.DB, bool) {
	if r == nil {
		return nil, false
	}
	db, ok := r.dbs[key]
	return db, ok
}

func (r *DatabaseRegistry) Default() (*gorm.DB, bool) {
	return r.Get(DEFAULT_KEY)
}

func (r *DatabaseRegistry) MustGet(key string) *gorm.DB {
	db, ok := r.Get(key)
	if !ok {
		panic(fmt.Sprintf("database %q not found", key))
	}
	return db
}

func (r *DatabaseRegistry) All() map[string]*gorm.DB {
	if r == nil {
		return nil
	}
	return maps.Clone(r.dbs)
}

func (r *DatabaseRegistry) Keys() []string {
	if r == nil {
		return nil
	}
	keys := maps.Keys(r.dbs)
	return slices.Collect(keys)
}

// CacheRegistry 提供对缓存实例集合的只读访问。
type CacheRegistry struct {
	caches map[string]cache.AdapterCache
}

func NewCacheRegistry(caches map[string]cache.AdapterCache) *CacheRegistry {
	return &CacheRegistry{
		caches: caches,
	}
}

func (r *CacheRegistry) Get(key string) (cache.AdapterCache, bool) {
	if r == nil {
		return nil, false
	}
	item, ok := r.caches[key]
	return item, ok
}

func (r *CacheRegistry) Default() (cache.AdapterCache, bool) {
	return r.Get(DEFAULT_KEY)
}

func (r *CacheRegistry) MustGet(key string) cache.AdapterCache {
	item, ok := r.Get(key)
	if !ok {
		panic(fmt.Sprintf("cache %q not found", key))
	}
	return item
}

func (r *CacheRegistry) All() map[string]cache.AdapterCache {
	if r == nil {
		return nil
	}
	return maps.Clone(r.caches)
}

func (r *CacheRegistry) Keys() []string {
	if r == nil {
		return nil
	}
	keys := maps.Keys(r.caches)
	return slices.Collect(keys)
}

// ConsumerRegistry 提供对 MQ consumer 集合的只读访问。
// 这里保留 map[string][]Consumer 的结构，是因为同一个 key 下本来就可能挂多个 consumer。
type ConsumerRegistry struct {
	consumers map[string][]rocketmq.Consumer
}

func NewConsumerRegistry(consumers map[string][]rocketmq.Consumer) *ConsumerRegistry {
	return &ConsumerRegistry{
		consumers: consumers,
	}
}

func (r *ConsumerRegistry) Get(key string) ([]rocketmq.Consumer, bool) {
	if r == nil {
		return nil, false
	}
	items, ok := r.consumers[key]
	if !ok {
		return nil, false
	}
	return append([]rocketmq.Consumer(nil), items...), true
}

func (r *ConsumerRegistry) All() map[string][]rocketmq.Consumer {
	if r == nil {
		return nil
	}
	clone := make(map[string][]rocketmq.Consumer, len(r.consumers))
	for key, items := range r.consumers {
		clone[key] = append([]rocketmq.Consumer(nil), items...)
	}
	return clone
}

func (r *ConsumerRegistry) Keys() []string {
	if r == nil {
		return nil
	}
	keys := maps.Keys(r.consumers)
	return slices.Collect(keys)
}
