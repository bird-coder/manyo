/*
 * @Author: yujiajie
 * @Date: 2024-05-23 17:39:39
 * @LastEditors: yujiajie
 * @LastEditTime: 2024-05-31 14:41:57
 * @FilePath: /manyo/pkg/pool/gopool/pool.go
 * @Description:
 */
package gopool

import "github.com/bird-coder/manyo/pkg/logger"

type Pool struct {
	ch       chan struct{}
	poolSize int
}

func NewPool(poolSize int) *Pool {
	pool := &Pool{
		ch:       make(chan struct{}, poolSize),
		poolSize: poolSize,
	}
	return pool
}

func (pool *Pool) Run(fn func()) {
	pool.ch <- struct{}{}
	go func() {
		defer func() {
			<-pool.ch
			if err := recover(); err != nil {
				logger.Error("goroutine run err: %v", err)
			}
		}()
		fn()
	}()
}
