/*
 * @Author: yujiajie
 * @Date: 2025-01-06 11:58:23
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-03-21 17:03:23
 * @FilePath: /manyo/pkg/zrpc/internal/clientinterceptors/sentinelinterceptor.go
 * @Description:
 */
package clientinterceptors

import (
	"context"
	"errors"
	"fmt"
	"path"

	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	cds "github.com/bird-coder/manyo/pkg/zrpc/internal/codes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryBreakerInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
	breakerName := path.Join(cc.Target(), method)
	if _, err := circuitbreaker.LoadRulesOfResource(breakerName, []*circuitbreaker.Rule{
		{
			Resource:         breakerName,
			Strategy:         circuitbreaker.ErrorRatio,
			RetryTimeoutMs:   3000,
			MinRequestAmount: 10,
			StatIntervalMs:   10000,
			Threshold:        0.3,
		},
	}); err != nil {
		fmt.Println(err)
	}
	entry, err := api.Entry(
		breakerName,
		api.WithResourceType(base.ResTypeRPC),
		api.WithTrafficType(base.Inbound),
	)
	if err != nil {
		err = status.Error(codes.Unavailable, err.Error())
		return
	}
	defer entry.Exit()

	err = invoker(ctx, method, req, reply, cc, opts...)

	if errors.Is(err, context.DeadlineExceeded) {
		err = status.Error(codes.Unavailable, err.Error())
		api.TraceError(entry, err)
	} else {
		if !cds.Acceptable(err) {
			api.TraceError(entry, err)
		}
	}
	return
}
