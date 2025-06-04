/*
 * @Author: yujiajie
 * @Date: 2025-06-04 11:47:04
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-06-04 11:47:19
 * @FilePath: /Go-Base/pkg/generator/local.go
 * @Description:
 */
package generator

import "sync"

type LocalGenerator struct {
	max     uint64
	current uint64
	used    map[uint64]bool

	lock sync.RWMutex
}

func NewGenerator(max uint64) *LocalGenerator {
	return &LocalGenerator{
		used:    make(map[uint64]bool, max),
		max:     max,
		current: 1,
	}
}

func (g *LocalGenerator) Next() uint64 {
	g.lock.Lock()
	defer g.lock.Unlock()

	for i := uint64(0); i < g.max; i++ {
		id := g.current
		if g.current == g.max {
			g.current = 1
		} else {
			g.current++
		}
		if !g.used[id] {
			g.used[id] = true
			return id
		}
	}

	return 0
}

func (g *LocalGenerator) Release(id uint64) {
	g.lock.Lock()
	defer g.lock.Unlock()

	delete(g.used, id)
}
