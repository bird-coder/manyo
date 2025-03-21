/*
 * @Author: yujiajie
 * @Date: 2025-03-21 17:01:02
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-03-21 17:01:36
 * @FilePath: /manyo/pkg/zrpc/proxy.go
 * @Description:
 */
package zrpc

import (
	"context"
	"sync"

	"github.com/bird-coder/manyo/pkg/syncx"
	"github.com/bird-coder/manyo/pkg/zrpc/internal"
	"google.golang.org/grpc"
)

type RpcProxy struct {
	backend      string
	clients      sync.Map
	options      []internal.ClientOption
	singleFlight syncx.SingleFlight
}

func NewProxy(backend string, opts ...internal.ClientOption) *RpcProxy {
	return &RpcProxy{
		backend:      backend,
		options:      opts,
		singleFlight: syncx.NewSingleFlight(),
	}
}

func (p *RpcProxy) TakeConn(ctx context.Context) (*grpc.ClientConn, error) {
	key := ""
	val, err := p.singleFlight.Do(key, func() (any, error) {
		client, ok := p.clients.Load(key)
		if ok {
			return client, nil
		}
		cli, err := NewClientWithTarget(p.backend, p.options...)
		if err != nil {
			return nil, err
		}
		p.clients.Store(key, cli)
		return cli, nil
	})
	if err != nil {
		return nil, err
	}
	return val.(internal.Client).Conn(), nil
}
