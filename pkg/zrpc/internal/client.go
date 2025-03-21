/*
 * @Author: yujiajie
 * @Date: 2025-01-02 17:48:55
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-03-21 17:02:45
 * @FilePath: /manyo/pkg/zrpc/internal/client.go
 * @Description:
 */
package internal

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bird-coder/manyo/pkg/zrpc/internal/balancer/p2c"
	"github.com/bird-coder/manyo/pkg/zrpc/internal/clientinterceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	dialTimeout = time.Second * 3
	separator   = '/'
)

type (
	Client interface {
		Conn() *grpc.ClientConn
	}

	ClientOptions struct {
		NonBlock    bool
		Timeout     time.Duration
		Secure      bool
		DialOptions []grpc.DialOption
	}

	ClientOption func(options *ClientOptions)

	client struct {
		conn        *grpc.ClientConn
		middlewares ClientMiddlewaresConf
	}
)

func NewClient(target string, middlewares ClientMiddlewaresConf, opts ...ClientOption) (Client, error) {
	cli := &client{
		middlewares: middlewares,
	}

	svcCfg := fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, p2c.Name)
	balancerOpt := WithDialOption(grpc.WithDefaultServiceConfig(svcCfg))
	opts = append([]ClientOption{balancerOpt}, opts...)
	if err := cli.dial(target, opts...); err != nil {
		return nil, err
	}

	return cli, nil
}

func (c *client) Conn() *grpc.ClientConn {
	return c.conn
}

func (c *client) buildDialOptions(opts ...ClientOption) []grpc.DialOption {
	var cliOpts ClientOptions
	for _, opt := range opts {
		opt(&cliOpts)
	}

	var options []grpc.DialOption
	if !cliOpts.Secure {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	options = append(options,
		grpc.WithChainUnaryInterceptor(c.buildUnaryInterceptors(cliOpts.Timeout)...),
		grpc.WithChainStreamInterceptor(c.buildStreamInterceptors()...))

	return append(options, cliOpts.DialOptions...)
}

func (c *client) dial(server string, opts ...ClientOption) error {
	options := c.buildDialOptions(opts...)
	conn, err := grpc.NewClient(server, options...)
	if err != nil {
		service := server
		if errors.Is(err, context.DeadlineExceeded) {
			pos := strings.LastIndexByte(server, separator)
			if pos > 0 && pos < len(server)-1 {
				service = server[pos+1:]
			}
		}
		return fmt.Errorf("rpc dial: %s, error: %s, make sure rpc service %q is already started",
			server, err.Error(), service)
	}

	c.conn = conn
	return nil
}

func (c *client) buildStreamInterceptors() []grpc.StreamClientInterceptor {
	var interceptors []grpc.StreamClientInterceptor

	return interceptors
}

func (c *client) buildUnaryInterceptors(timeout time.Duration) []grpc.UnaryClientInterceptor {
	var interceptors []grpc.UnaryClientInterceptor

	if c.middlewares.Duration {
		interceptors = append(interceptors, clientinterceptors.UnaryDurationInterceptor)
	}
	if c.middlewares.Prometheus {
		interceptors = append(interceptors, clientinterceptors.UnaryPrometheusInterceptor)
	}
	if c.middlewares.Breaker {
		interceptors = append(interceptors, clientinterceptors.UnaryBreakerInterceptor)
	}
	if c.middlewares.Timeout {
		interceptors = append(interceptors, clientinterceptors.UnaryTimeoutInterceptor(timeout))
	}
	return interceptors
}

func WithDialOption(opt grpc.DialOption) ClientOption {
	return func(options *ClientOptions) {
		options.DialOptions = append(options.DialOptions, opt)
	}
}

func WithNonBlock() ClientOption {
	return func(options *ClientOptions) {
		options.NonBlock = true
	}
}

func WithStreamClientInterceptor(interceptor grpc.StreamClientInterceptor) ClientOption {
	return func(options *ClientOptions) {
		options.DialOptions = append(options.DialOptions, grpc.WithChainStreamInterceptor(interceptor))
	}
}

func WithUnaryClientInterceptor(interceptor grpc.UnaryClientInterceptor) ClientOption {
	return func(options *ClientOptions) {
		options.DialOptions = append(options.DialOptions, grpc.WithChainUnaryInterceptor(interceptor))
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(options *ClientOptions) {
		options.Timeout = timeout
	}
}

func WithTransportCredentials(creds credentials.TransportCredentials) ClientOption {
	return func(options *ClientOptions) {
		options.Secure = true
		options.DialOptions = append(options.DialOptions, grpc.WithTransportCredentials(creds))
	}
}
