/*
 * @Author: yujiajie
 * @Date: 2025-01-03 14:12:21
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-03-21 17:02:09
 * @FilePath: /manyo/pkg/zrpc/internal/serverinterceptors/sentinelinterceptor.go
 * @Description:
 */
package serverinterceptors

import (
	"context"
	"errors"
	"fmt"

	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	cds "github.com/bird-coder/manyo/pkg/zrpc/internal/codes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func StreamBreakerInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	breakerName := info.FullMethod
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

	err = handler(srv, ss)

	if !cds.Acceptable(err) {
		api.TraceError(entry, err)
	}
	return
}

func UnaryBreakerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	breakerName := info.FullMethod
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

	resp, err = handler(ctx, req)

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
