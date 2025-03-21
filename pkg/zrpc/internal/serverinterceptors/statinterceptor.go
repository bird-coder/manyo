/*
 * @Author: yujiajie
 * @Date: 2025-01-03 15:13:29
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-03-21 17:02:20
 * @FilePath: /manyo/pkg/zrpc/internal/serverinterceptors/statinterceptor.go
 * @Description:
 */
package serverinterceptors

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/bird-coder/manyo/pkg/collection"
	"github.com/bird-coder/manyo/pkg/logger"
	"github.com/bird-coder/manyo/pkg/syncx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

const defaultSlowThreshold = time.Millisecond * 500

var (
	ignoreContentMethods sync.Map
	slowThreshold        = syncx.ForAtomicDuration(defaultSlowThreshold)
)

type StatConf struct {
	SlowThreshold        time.Duration `json:",default=500ms"`
	IgnoreContentMethods []string      `json:",optional"`
}

func DontLogContentForMethod(method string) {
	ignoreContentMethods.Store(method, struct{}{})
}

func SetSlowThreshold(threshold time.Duration) {
	slowThreshold.Set(threshold)
}

func UnaryStatInterceptor(conf StatConf) grpc.UnaryServerInterceptor {
	notLogMethods := collection.NewSet()
	notLogMethods.AddStr(conf.IgnoreContentMethods...)
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		startTime := time.Now()
		defer func() {
			duration := time.Since(startTime)
			logDuration(ctx, info.FullMethod, req, duration, notLogMethods, conf.SlowThreshold)
		}()

		return handler(ctx, req)
	}
}

func isSlow(duration, durationThreshold time.Duration) bool {
	return duration > slowThreshold.Load() ||
		(durationThreshold > 0 && duration > durationThreshold)
}

func logDuration(ctx context.Context, method string, req any, duration time.Duration,
	ignoreMethods *collection.Set, durationThreshold time.Duration) {
	var addr string
	client, ok := peer.FromContext(ctx)
	if ok {
		addr = client.Addr.String()
	}
	if !shouldLogContent(method, ignoreMethods) {
		if isSlow(duration, durationThreshold) {
			logger.Error("[RPC] slowcall - %s - %s", addr, method)
		}
	} else {
		content, err := json.Marshal(req)
		if err != nil {
			logger.Error("%s - %s", addr, err.Error())
		} else if duration > defaultSlowThreshold {
			logger.Error("[RPC] slowcall - %s - %s - %s", addr, method, string(content))
		} else {
			logger.Info("%s - %s - %s", addr, method, string(content))
		}
	}
}

func shouldLogContent(method string, ignoreMethods *collection.Set) bool {
	_, ok := ignoreContentMethods.Load(method)
	return !ok && !ignoreMethods.Contains(method)
}
