/*
 * @Author: yujiajie
 * @Date: 2025-01-02 17:19:05
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-01-02 17:48:29
 * @FilePath: /Go-Base/pkg/zrpc/internal/rpcserver.go
 * @Description:
 */
package internal

import (
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const defaultConnectionIdleDuration = time.Minute * 5

type (
	RegisterFn func(*grpc.Server)

	Server interface {
		AddOptions(options ...grpc.ServerOption)
		AddStreamInterceptors(interceptors ...grpc.StreamServerInterceptor)
		AddUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor)
		SetName(string)
		Start(register RegisterFn) error
	}

	rpcServer struct {
		name               string
		address            string
		options            []grpc.ServerOption
		streamInterceptors []grpc.StreamServerInterceptor
		unaryInterceptors  []grpc.UnaryServerInterceptor
	}
)

func NewRpcServer(addr string) Server {
	return &rpcServer{
		address: addr,
		options: []grpc.ServerOption{grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: defaultConnectionIdleDuration,
		})},
	}
}

func (s *rpcServer) SetName(name string) {
	s.name = name
}

func (s *rpcServer) Start(register RegisterFn) error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	unaryInterceptorOption := grpc.ChainUnaryInterceptor(s.unaryInterceptors...)
	streamInterceptorOption := grpc.ChainStreamInterceptor(s.streamInterceptors...)

	options := append(s.options, unaryInterceptorOption, streamInterceptorOption)
	server := grpc.NewServer(options...)
	register(server)

	defer server.GracefulStop()

	return server.Serve(lis)
}

func (s *rpcServer) AddOptions(options ...grpc.ServerOption) {
	s.options = append(s.options, options...)
}

func (s *rpcServer) AddStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) {
	s.streamInterceptors = append(s.streamInterceptors, interceptors...)
}

func (s *rpcServer) AddUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) {
	s.unaryInterceptors = append(s.unaryInterceptors, interceptors...)
}
