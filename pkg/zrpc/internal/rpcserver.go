/*
 * @Author: yujiajie
 * @Date: 2025-01-02 17:19:05
 * @LastEditors: yujiajie 1037297660@qq.com
 * @LastEditTime: 2026-05-09 10:44:26
 * @FilePath: /manyo/pkg/zrpc/internal/rpcserver.go
 * @Description:
 */
package internal

import (
	"context"
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
		Prepare(context.Context) error
		Start(context.Context, RegisterFn) error
		Stop(context.Context) error
	}

	rpcServer struct {
		name               string
		address            string
		options            []grpc.ServerOption
		streamInterceptors []grpc.StreamServerInterceptor
		unaryInterceptors  []grpc.UnaryServerInterceptor

		server *grpc.Server
		ln     net.Listener
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

func (s *rpcServer) Prepare(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	unaryInterceptorOption := grpc.ChainUnaryInterceptor(s.unaryInterceptors...)
	streamInterceptorOption := grpc.ChainStreamInterceptor(s.streamInterceptors...)

	options := append(s.options, unaryInterceptorOption, streamInterceptorOption)
	server := grpc.NewServer(options...)

	s.server = server
	s.ln = ln

	return nil
}

func (s *rpcServer) Start(ctx context.Context, register RegisterFn) error {
	register(s.server)

	return s.server.Serve(s.ln)
}

func (s *rpcServer) Stop(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		s.server.GracefulStop()
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		s.server.Stop()
		<-done
		return ctx.Err()
	}
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
