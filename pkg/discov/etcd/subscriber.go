/*
 * @Author: yujiajie
 * @Date: 2025-01-23 17:40:46
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-03-21 16:59:56
 * @FilePath: /manyo/pkg/discov/etcd/subscriber.go
 * @Description:
 */
package etcd

import (
	"sync"
	"sync/atomic"

	"github.com/bird-coder/manyo/pkg/discov/etcd/internal"
	"github.com/bird-coder/manyo/pkg/logger"
	"github.com/bird-coder/manyo/pkg/syncx"
)

type (
	SubOption func(sub *SubScriber)

	SubScriber struct {
		endpoints  []string
		exclusive  bool
		exactMatch bool
		items      *container
	}
)

func NewSubscriber(endpoints []string, key string, opts ...SubOption) (*SubScriber, error) {
	sub := &SubScriber{
		endpoints: endpoints,
	}
	for _, opt := range opts {
		opt(sub)
	}
	sub.items = newContainer(sub.exclusive)

	if err := internal.GetRegistry().Monitor(endpoints, key, sub.items, sub.exactMatch); err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *SubScriber) AddListener(listener func()) {
	s.items.addListener(listener)
}

func (s *SubScriber) Values() []string {
	return s.items.getValues()
}

func Exclusive() SubOption {
	return func(sub *SubScriber) {
		sub.exclusive = true
	}
}

func WithSubEtcdAccount(user, pass string) SubOption {
	return func(sub *SubScriber) {
		RegisterAccount(sub.endpoints, user, pass)
	}
}

func WithSubEtcdTLS(certFile, certKeyFile, caFile string, insecureSkipVerify bool) SubOption {
	return func(sub *SubScriber) {
		err := RegisterTLS(sub.endpoints, certFile, certKeyFile, caFile, insecureSkipVerify)
		logger.Error("%v", err)
	}
}

type container struct {
	exclusive bool
	values    map[string][]string
	mapping   map[string]string
	snapshot  atomic.Value
	dirty     *syncx.AtomicBool
	listeners []func()
	lock      sync.Mutex
}

func newContainer(exclusive bool) *container {
	return &container{
		exclusive: exclusive,
		values:    make(map[string][]string),
		mapping:   make(map[string]string),
		dirty:     syncx.ForAtomicBool(true),
	}
}

func (c *container) OnAdd(kv internal.KV) {
	c.addKv(kv.Key, kv.Val)
	c.notifyChange()
}

func (c *container) OnDelete(kv internal.KV) {
	c.removeKey(kv.Key)
	c.notifyChange()
}

func (c *container) addKv(key, value string) ([]string, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.dirty.Set(true)
	keys := c.values[value]
	previous := append([]string(nil), keys...)
	early := len(keys) > 0
	if c.exclusive && early {
		for _, each := range keys {
			c.doRemoveKey(each)
		}
	}
	c.values[value] = append(c.values[value], key)
	c.mapping[key] = value

	if early {
		return previous, true
	}
	return nil, false
}

func (c *container) removeKey(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.dirty.Set(true)
	c.doRemoveKey(key)
}

func (c *container) addListener(listener func()) {
	c.lock.Lock()
	c.listeners = append(c.listeners, listener)
	c.lock.Unlock()
}

func (c *container) doRemoveKey(key string) {
	server, ok := c.mapping[key]
	if !ok {
		return
	}

	delete(c.mapping, key)
	keys := c.values[server]
	remain := keys[:0]

	for _, k := range keys {
		if k != key {
			remain = append(remain, k)
		}
	}

	if len(remain) > 0 {
		c.values[server] = remain
	} else {
		delete(c.values, server)
	}
}

func (c *container) getValues() []string {
	if !c.dirty.True() {
		return c.snapshot.Load().([]string)
	}
	c.lock.Lock()
	defer c.lock.Unlock()

	var vals []string
	for each := range c.values {
		vals = append(vals, each)
	}
	c.snapshot.Store(vals)
	c.dirty.Set(false)

	return vals
}

func (c *container) notifyChange() {
	c.lock.Lock()
	listeners := append(([]func())(nil), c.listeners...)
	c.lock.Unlock()

	for _, listener := range listeners {
		listener()
	}
}
