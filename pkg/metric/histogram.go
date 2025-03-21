/*
 * @Description:
 * @Author: yujiajie
 * @Date: 2023-11-30 00:01:59
 * @LastEditTime: 2024-03-26 11:22:06
 * @LastEditors: yujiajie
 */
package metric

import (
	prom "github.com/prometheus/client_golang/prometheus"
)

type HistogramVecOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
	Buckets   []float64
}

type HistogramVec interface {
	Observe(v int64, labels ...string)
	close() bool
}

type promHistogramVec struct {
	histogram *prom.HistogramVec
}

// 创建prometheus统计指标
func NewHistogramVec(cfg *HistogramVecOpts) HistogramVec {
	if cfg == nil {
		return nil
	}
	vec := prom.NewHistogramVec(prom.HistogramOpts{
		Namespace: cfg.Namespace,
		Subsystem: cfg.Subsystem,
		Name:      cfg.Name,
		Help:      cfg.Help,
		Buckets:   cfg.Buckets,
	}, cfg.Labels)
	prom.MustRegister(vec)
	hv := &promHistogramVec{
		histogram: vec,
	}
	return hv
}

func (hv *promHistogramVec) Observe(v int64, labels ...string) {
	hv.histogram.WithLabelValues(labels...).Observe(float64(v))
}

func (hv *promHistogramVec) close() bool {
	return prom.Unregister(hv.histogram)
}
