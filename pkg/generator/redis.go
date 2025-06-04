/*
 * @Author: yujiajie
 * @Date: 2025-06-03 18:19:26
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-06-04 12:03:28
 * @FilePath: /Go-Base/pkg/generator/redis.go
 * @Description:
 */
package generator

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	luaAcquire = redis.NewScript(`
		local max = tonumber(ARGV[1])
		local tag = tonumber(ARGV[2])
		local ttl = tonumber(ARGV[3])

		for i = 1, max do
			local r1 = redis.call("incr", KEYS[1])
			r1 = ((r1-1) % max) + 1

			local key = KEYS[2]..r1
			if redis.call("set", key, tag, "NX", "EX", ttl) then
				return r1
			end
		end
		return redis.error_reply("ID generation failed")
	`)
	luaTTL = redis.NewScript(`
		local ttl = tonumber(ARGV[1])
		for i, key in ipairs(KEYS) do
			redis.call("EXPIRE", key, ttl)
		end
		return redis.status_reply("OK")
	`)
)

var (
	appName   = "rock"
	serverKey = "nbgame:server_id:%s"
	activeKey = "nbgame:server_id:%s:active:"
)

type options struct {
	ctx         context.Context
	appName     string
	serverKey   string
	activeKey   string
	expireTime  int
	refreshTime int
}

type Option func(o *options)

func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

func WithAppName(appName string) Option {
	return func(o *options) {
		o.appName = appName
	}
}

func WithServerKey(serverKey string) Option {
	return func(o *options) {
		o.serverKey = serverKey
	}
}

func WithActiveKey(activeKey string) Option {
	return func(o *options) {
		o.activeKey = activeKey
	}
}

func WithExpireTime(expireTime int) Option {
	return func(o *options) {
		o.expireTime = expireTime
	}
}

func WithRefreshTime(refreshTime int) Option {
	return func(o *options) {
		o.refreshTime = refreshTime
	}
}

type RedisGenerator struct {
	ctx    context.Context
	cancel context.CancelFunc

	opts options

	rds   *redis.Client
	maxId uint64

	used map[uint64]bool

	tick *time.Ticker

	lock sync.RWMutex
}

func NewRedisGenerator(rds *redis.Client, maxId uint64, opts ...Option) *RedisGenerator {
	if rds == nil {
		return nil
	}
	op := options{
		ctx:         context.Background(),
		appName:     appName,
		serverKey:   fmt.Sprintf(serverKey, appName),
		activeKey:   fmt.Sprintf(activeKey, appName),
		expireTime:  15,
		refreshTime: 10,
	}
	for _, opt := range opts {
		opt(&op)
	}
	if maxId > math.MaxInt64 {
		maxId = math.MaxInt64
	}
	ctx, cancel := context.WithCancel(op.ctx)
	rg := &RedisGenerator{
		ctx:    ctx,
		cancel: cancel,
		opts:   op,
		rds:    rds,
		maxId:  maxId,
		used:   make(map[uint64]bool),
		tick:   time.NewTicker(time.Duration(op.refreshTime) * time.Second),
	}
	go rg.keepAlive()

	return rg
}

func (g *RedisGenerator) Next() uint64 {
	id, err := g.acquireID()
	if err != nil {
		return 0
	}
	g.lock.Lock()
	g.used[id] = true
	g.lock.Unlock()

	return id
}

func (g *RedisGenerator) Release(id uint64) {
	g.lock.Lock()
	delete(g.used, id)
	g.lock.Unlock()

	g.rds.Del(g.ctx, g.formatActiveKey(id))
}

func (g *RedisGenerator) formatActiveKey(id uint64) string {
	return g.opts.activeKey + strconv.FormatInt(int64(id), 10)
}

func (g *RedisGenerator) acquireID() (uint64, error) {
	res, err := luaAcquire.Run(g.ctx, g.rds, []string{g.opts.serverKey, g.opts.activeKey}, g.maxId, 1, g.opts.expireTime).Result()
	if err != nil {
		return 0, err
	}
	id := res.(int64)
	return uint64(id), nil
}

func (g *RedisGenerator) keepAlive() {
	defer g.tick.Stop()
	for {
		select {
		case <-g.tick.C:
			keys := g.getUsedKey()
			if len(keys) > 0 {
				luaTTL.Run(g.ctx, g.rds, keys, g.opts.expireTime)
			}
		case <-g.ctx.Done():
			return
		}
	}
}

func (rg *RedisGenerator) getUsedKey() []string {
	rg.lock.RLock()
	defer rg.lock.RUnlock()
	keys := make([]string, 0, len(rg.used))
	for k, v := range rg.used {
		if v {
			keys = append(keys, rg.formatActiveKey(k))
		}
	}
	return keys
}
