/*
 * @Author: yujiajie
 * @Date: 2025-07-11 17:26:34
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-07-24 09:39:42
 * @FilePath: /Go-Base/lib/stage/options.go
 * @Description:
 */
package stage

import (
	"context"
	"time"
)

type Server interface {
	Start() error
	Stop() error
}

type optionFunc func(o *options)

type options struct {
	ctx context.Context

	stopTimeout time.Duration

	servers []Server

	beforeStart []func() error
	afterStart  []func() error
	beforeStop  []func() error
	afterStop   []func() error
}

func WithContext(ctx context.Context) optionFunc {
	return func(o *options) {
		o.ctx = ctx
	}
}

func WithStopTimeout(timeout time.Duration) optionFunc {
	return func(o *options) {
		o.stopTimeout = timeout
	}
}

func WithServer(servers ...Server) optionFunc {
	return func(o *options) {
		o.servers = servers
	}
}

func BeforeStart(fn func() error) optionFunc {
	return func(o *options) {
		o.beforeStart = append(o.beforeStart, fn)
	}
}

func AfterStart(fn func() error) optionFunc {
	return func(o *options) {
		o.afterStart = append(o.afterStart, fn)
	}
}

func BeforeStop(fn func() error) optionFunc {
	return func(o *options) {
		o.beforeStop = append(o.beforeStop, fn)
	}
}

func AfterStop(fn func() error) optionFunc {
	return func(o *options) {
		o.afterStop = append(o.afterStop, fn)
	}
}
