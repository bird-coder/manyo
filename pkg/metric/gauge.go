/*
 * @Author: yujiajie
 * @Date: 2025-06-18 10:33:19
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-06-18 10:41:02
 * @FilePath: /Go-Base/pkg/metric/gauge.go
 * @Description:
 */
package metric

import prom "github.com/prometheus/client_golang/prometheus"

type GaugeVectorOpts struct {
	Namespace string
	Subsystem string
	Name      string
	Help      string
	Labels    []string
}

type GaugeVec interface {
	Inc(labels ...string)
	Add(v float64, labels ...string)
	Sub(v float64, labels ...string)
	Set(v float64, labels ...string)
	close() bool
}

type promGaugeVec struct {
	gauge *prom.GaugeVec
}

// 创建prometheus统计指标
func NewGaugeVec(cfg *GaugeVectorOpts) GaugeVec {
	if cfg == nil {
		return nil
	}
	vec := prom.NewGaugeVec(prom.GaugeOpts{
		Namespace: cfg.Namespace,
		Subsystem: cfg.Subsystem,
		Name:      cfg.Name,
		Help:      cfg.Help,
	}, cfg.Labels)
	prom.MustRegister(vec)
	cv := &promGaugeVec{
		gauge: vec,
	}

	return cv
}

func (cv *promGaugeVec) Inc(labels ...string) {
	cv.gauge.WithLabelValues(labels...).Inc()
}

func (cv *promGaugeVec) Add(v float64, labels ...string) {
	cv.gauge.WithLabelValues(labels...).Add(v)
}

func (cv *promGaugeVec) Sub(v float64, labels ...string) {
	cv.gauge.WithLabelValues(labels...).Sub(v)
}

func (cv *promGaugeVec) Set(v float64, labels ...string) {
	cv.gauge.WithLabelValues(labels...).Set(v)
}

func (cv *promGaugeVec) close() bool {
	return prom.Unregister(cv.gauge)
}
