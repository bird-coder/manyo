/*
 * @Author: yujiajie
 * @Date: 2025-01-06 11:42:39
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-03-21 17:03:05
 * @FilePath: /manyo/pkg/zrpc/internal/clientinterceptors/durationinterceptor.go
 * @Description:
 */
package clientinterceptors

import (
	"context"
	"path"
	"sync"
	"time"

	"github.com/bird-coder/manyo/pkg/logger"
	"github.com/bird-coder/manyo/pkg/syncx"
	"google.golang.org/grpc"
)

const defaultSlowThreshold = time.Millisecond * 500

var (
	ignoreContentMethods sync.Map
	slowThreshold        = syncx.ForAtomicDuration(defaultSlowThreshold)
)

func UnaryDurationInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	serverName := path.Join(cc.Target(), method)
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	duration := time.Since(start)
	_, ok := ignoreContentMethods.Load(method)
	if err != nil {
		if ok {
			logger.Error("fail - %s - %s", serverName, err.Error())
		} else {
			logger.Error("fail - %s - %v - %s", serverName, req, err.Error())
		}
	} else if duration > slowThreshold.Load() {
		if ok {
			logger.Error("[RPC] ok - slowcall - %s", serverName)
		} else {
			logger.Error("[RPC] ok - slowcall - %s - %v - %v", serverName, req, reply)
		}
	}
	return err
}

func DontLogContentForMethod(method string) {
	ignoreContentMethods.Store(method, struct{}{})
}

func SetSlowThreshold(threshold time.Duration) {
	slowThreshold.Set(threshold)
}
