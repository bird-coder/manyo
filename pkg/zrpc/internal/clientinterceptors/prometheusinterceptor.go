/*
 * @Author: yujiajie
 * @Date: 2025-01-06 11:33:30
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-03-21 17:03:17
 * @FilePath: /manyo/pkg/zrpc/internal/clientinterceptors/prometheusinterceptor.go
 * @Description:
 */
package clientinterceptors

import (
	"context"
	"strconv"
	"time"

	"github.com/bird-coder/manyo/pkg/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const clientNamespace = "rpc_client"

var (
	metricClientReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: clientNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "rpc client requests duration(ms).",
		Labels:    []string{"method"},
		Buckets:   []float64{1, 2, 5, 10, 25, 50, 100, 250, 500, 1000, 2000, 5000},
	})

	metricClientReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: clientNamespace,
		Subsystem: "requests",
		Name:      "code_total",
		Help:      "rpc client requests code count.",
		Labels:    []string{"method", "code"},
	})
)

func UnaryPrometheusInterceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	startTime := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	metricClientReqDur.Observe(time.Since(startTime).Milliseconds(), method)
	metricClientReqCodeTotal.Inc(method, strconv.Itoa(int(status.Code(err))))
	return err
}
