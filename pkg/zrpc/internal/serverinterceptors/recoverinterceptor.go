/*
 * @Author: yujiajie
 * @Date: 2025-01-03 10:31:03
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-03-21 17:02:03
 * @FilePath: /manyo/pkg/zrpc/internal/serverinterceptors/recoverinterceptor.go
 * @Description:
 */
package serverinterceptors

import (
	"context"
	"runtime/debug"

	"github.com/bird-coder/manyo/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func StreamRecoverInterceptor(svr any, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	defer handleCrash(func(r any) {
		err = toPanicError(r)
	})
	return handler(svr, stream)
}

func UnaryRecoverInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer handleCrash(func(r any) {
		err = toPanicError(r)
	})
	return handler(ctx, req)
}

func handleCrash(handler func(any)) {
	if r := recover(); r != nil {
		handler(r)
	}
}

func toPanicError(r any) error {
	logger.Error("%+v\n\n%s", r, debug.Stack())
	return status.Errorf(codes.Internal, "panic: %v", r)
}
