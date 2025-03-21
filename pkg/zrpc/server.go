/*
 * @Author: yujiajie
 * @Date: 2025-01-02 16:58:01
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-03-21 17:01:21
 * @FilePath: /manyo/pkg/zrpc/server.go
 * @Description:
 */
package zrpc

import (
	"time"

	"github.com/bird-coder/manyo/pkg/zrpc/internal"
	"github.com/bird-coder/manyo/pkg/zrpc/internal/serverinterceptors"
)

type RpcServer struct {
	server   internal.Server
	register internal.RegisterFn
}

func NewServer(c RpcServerConf, register internal.RegisterFn) (*RpcServer, error) {
	var server internal.Server

	server = internal.NewRpcServer(c.Addr)

	server.SetName(c.Name)
	setupStreamInterceptors(server, c)
	setupUnaryInterceptors(server, c)

	rpcServer := &RpcServer{
		server:   server,
		register: register,
	}

	return rpcServer, nil
}

func (rs *RpcServer) Start() {
	rs.server.Start(rs.register)
}

func (rs *RpcServer) Stop() {

}

func setupStreamInterceptors(svr internal.Server, c RpcServerConf) {
	if c.Middlewares.Recover {
		svr.AddStreamInterceptors(serverinterceptors.StreamRecoverInterceptor)
	}
	if c.Middlewares.Breaker {
		svr.AddStreamInterceptors(serverinterceptors.StreamBreakerInterceptor)
	}
}

func setupUnaryInterceptors(svr internal.Server, c RpcServerConf) {
	if c.Middlewares.Recover {
		svr.AddUnaryInterceptors(serverinterceptors.UnaryRecoverInterceptor)
	}
	if c.Middlewares.Stat {
		svr.AddUnaryInterceptors(serverinterceptors.UnaryStatInterceptor(c.Middlewares.StatConf))
	}
	if c.Middlewares.Prometheus {
		svr.AddUnaryInterceptors(serverinterceptors.UnaryPrometheusInterceptor)
	}
	if c.Middlewares.Breaker {
		svr.AddUnaryInterceptors(serverinterceptors.UnaryBreakerInterceptor)
	}
	if c.Timeout > 0 {
		svr.AddUnaryInterceptors(serverinterceptors.UnaryTimeoutInterceptor(
			time.Duration(c.Timeout)*time.Millisecond, c.MethodTimeouts...))
	}
}
