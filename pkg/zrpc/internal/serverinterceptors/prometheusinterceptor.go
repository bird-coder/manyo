/*
 * @Author: yujiajie
 * @Date: 2025-01-03 10:49:06
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-03-21 17:01:57
 * @FilePath: /manyo/pkg/zrpc/internal/serverinterceptors/prometheusinterceptor.go
 * @Description:
 */
package serverinterceptors

import (
	"context"
	"strconv"
	"time"

	"github.com/bird-coder/manyo/pkg/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const serverNamespace = "rpc_server"

var (
	metricServerReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "rpc server requests duration(ms).",
		Labels:    []string{"method"},
		Buckets:   []float64{1, 2, 5, 10, 25, 50, 100, 250, 500, 1000, 2000, 5000},
	})

	metricServerReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "rpc server requests code count.",
		Labels:    []string{"method", "code"},
	})
)

func UnaryPrometheusInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	startTime := time.Now()
	resp, err = handler(ctx, req)
	metricServerReqDur.Observe(time.Since(startTime).Milliseconds(), info.FullMethod)
	metricServerReqCodeTotal.Inc(info.FullMethod, strconv.Itoa(int(status.Code(err))))
	return resp, err
}
