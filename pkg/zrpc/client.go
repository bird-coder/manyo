/*
 * @Author: yujiajie
 * @Date: 2025-01-06 09:53:59
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-03-21 17:01:42
 * @FilePath: /manyo/pkg/zrpc/client.go
 * @Description:
 */
package zrpc

import (
	"time"

	"github.com/bird-coder/manyo/pkg/zrpc/internal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

type (
	RpcClient struct {
		client internal.Client
	}
)

func NewClient(c RpcClientConf, options ...internal.ClientOption) (internal.Client, error) {
	var opts []internal.ClientOption

	if c.NonBlock {
		opts = append(opts, internal.WithNonBlock())
	}
	if c.Timeout > 0 {
		opts = append(opts, internal.WithTimeout(time.Duration(c.Timeout)*time.Millisecond))
	}
	if c.KeepaliveTime > 0 {
		opts = append(opts, internal.WithDialOption(grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time: c.KeepaliveTime,
		})))
	}

	opts = append(opts, options...)

	client, err := internal.NewClient(c.Target, c.Middlewares, opts...)
	if err != nil {
		return nil, err
	}

	return &RpcClient{
		client: client,
	}, nil
}

func NewClientWithTarget(target string, opts ...internal.ClientOption) (internal.Client, error) {
	var config = RpcClientConf{
		Target: target,
	}

	return NewClient(config, opts...)
}

func (rc *RpcClient) Conn() *grpc.ClientConn {
	return rc.client.Conn()
}
