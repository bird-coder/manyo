/*
 * @Author: yujiajie
 * @Date: 2025-01-02 17:55:16
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-01-03 10:23:57
 * @FilePath: /Go-Base/pkg/zrpc/internal/serverinterceptors/timeoutinterceptor.go
 * @Description:
 */
package serverinterceptors

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	MethodTimeoutConf struct {
		FullMethod string
		Timeout    time.Duration
	}

	methodTimeouts map[string]time.Duration
)

func UnaryTimeoutInterceptor(timeout time.Duration, methodTimeouts ...MethodTimeoutConf) grpc.UnaryServerInterceptor {
	timeouts := buildMethodTimeouts(methodTimeouts)
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		t := getTimeoutByUnaryServerInfo(info.FullMethod, timeouts, timeout)
		ctx, cancel := context.WithTimeout(ctx, t)
		defer cancel()

		var lock sync.Mutex
		done := make(chan struct{})
		panicChan := make(chan any, 1)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- fmt.Sprintf("%+v\n\n%s", p, strings.TrimSpace(string(debug.Stack())))
				}
			}()

			lock.Lock()
			defer lock.Unlock()
			resp, err = handler(ctx, req)
			close(done)
		}()

		select {
		case p := <-panicChan:
			panic(p)
		case <-done:
			lock.Lock()
			defer lock.Unlock()
			return resp, err
		case <-ctx.Done():
			err := ctx.Err()
			if errors.Is(err, context.Canceled) {
				err = status.Error(codes.Canceled, err.Error())
			} else if errors.Is(err, context.DeadlineExceeded) {
				err = status.Error(codes.DeadlineExceeded, err.Error())
			}
			return nil, err
		}
	}
}

func buildMethodTimeouts(timeouts []MethodTimeoutConf) methodTimeouts {
	mt := make(methodTimeouts, len(timeouts))
	for _, st := range timeouts {
		if st.FullMethod != "" {
			mt[st.FullMethod] = st.Timeout
		}
	}
	return mt
}

func getTimeoutByUnaryServerInfo(method string, timeouts methodTimeouts,
	defaultTimeout time.Duration) time.Duration {
	if v, ok := timeouts[method]; ok {
		return v
	}

	return defaultTimeout
}
