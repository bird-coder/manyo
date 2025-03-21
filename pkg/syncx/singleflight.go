/*
 * @Author: yujiajie
 * @Date: 2025-01-23 16:07:17
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-01-23 16:15:20
 * @FilePath: /Go-Base/pkg/syncx/singleflight.go
 * @Description:
 */
package syncx

import "sync"

type (
	SingleFlight interface {
		Do(key string, fn func() (any, error)) (any, error)
		DoEx(key string, fn func() (any, error)) (any, bool, error)
	}

	call struct {
		wg  sync.WaitGroup
		val any
		err error
	}

	flightGroup struct {
		calls map[string]*call
		lock  sync.Mutex
	}
)

func NewSingleFlight() SingleFlight {
	return &flightGroup{
		calls: make(map[string]*call),
	}
}

func (g *flightGroup) Do(key string, fn func() (any, error)) (any, error) {
	c, done := g.createCall(key)
	if done {
		return c.val, c.err
	}

	g.makeCall(c, key, fn)
	return c.val, c.err
}

func (g *flightGroup) DoEx(key string, fn func() (any, error)) (val any, fresh bool, err error) {
	c, done := g.createCall(key)
	if done {
		return c.val, false, c.err
	}

	g.makeCall(c, key, fn)
	return c.val, true, c.err
}

func (g *flightGroup) createCall(key string) (c *call, done bool) {
	g.lock.Lock()
	if c, ok := g.calls[key]; ok {
		g.lock.Unlock()
		c.wg.Wait()
		return c, true
	}

	c = new(call)
	c.wg.Add(1)
	g.calls[key] = c
	g.lock.Unlock()

	return c, false
}

func (g *flightGroup) makeCall(c *call, key string, fn func() (any, error)) {
	defer func() {
		g.lock.Lock()
		delete(g.calls, key)
		g.lock.Unlock()
		c.wg.Done()
	}()

	c.val, c.err = fn()
}
