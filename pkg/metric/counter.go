/*
 * @Description:
 * @Author: yujiajie
 * @Date: 2023-12-05 21:46:52
 * @LastEditTime: 2024-03-26 11:22:02
 * @LastEditors: yujiajie
 */
package metric

import (
	prom "github.com/prometheus/client_golang/prometheus"
)

type VectorOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
}

type CounterVecOpts VectorOpts

type CounterVec interface {
	Inc(labels ...string)
	Add(v float64, labels ...string)
	close() bool
}

type promCounterVec struct {
	counter *prom.CounterVec
}

// 创建prometheus统计指标
func NewCounterVec(cfg *CounterVecOpts) CounterVec {
	if cfg == nil {
		return nil
	}
	vec := prom.NewCounterVec(prom.CounterOpts{
		Namespace: cfg.Namespace,
		Subsystem: cfg.Subsystem,
		Name:      cfg.Name,
		Help:      cfg.Help,
	}, cfg.Labels)
	prom.MustRegister(vec)
	cv := &promCounterVec{
		counter: vec,
	}

	return cv
}

func (cv *promCounterVec) Inc(labels ...string) {
	cv.counter.WithLabelValues(labels...).Inc()
}

func (cv *promCounterVec) Add(v float64, labels ...string) {
	cv.counter.WithLabelValues(labels...).Add(v)
}

func (cv *promCounterVec) close() bool {
	return prom.Unregister(cv.counter)
}
