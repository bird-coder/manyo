/*
 * @Author: yujiajie
 * @Date: 2025-01-06 11:38:44
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-01-06 11:42:10
 * @FilePath: /Go-Base/pkg/zrpc/internal/clientinterceptors/timeoutinterceptor.go
 * @Description:
 */
package clientinterceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

type TimeoutCallOption struct {
	grpc.EmptyCallOption
	timeout time.Duration
}

func UnaryTimeoutInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		t := getTimeoutFromCallOptions(opts, timeout)
		if t <= 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		ctx, cancel := context.WithTimeout(ctx, t)
		defer cancel()

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func WithCallTimeout(timeout time.Duration) grpc.CallOption {
	return TimeoutCallOption{
		timeout: timeout,
	}
}

func getTimeoutFromCallOptions(opts []grpc.CallOption, defaultTimeout time.Duration) time.Duration {
	for _, opt := range opts {
		if o, ok := opt.(TimeoutCallOption); ok {
			return o.timeout
		}
	}

	return defaultTimeout
}
