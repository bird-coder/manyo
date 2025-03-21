/*
 * @Author: yujiajie
 * @Date: 2025-03-21 17:01:02
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-03-21 17:02:52
 * @FilePath: /manyo/pkg/zrpc/internal/config.go
 * @Description:
 */
package internal

import "github.com/bird-coder/manyo/pkg/zrpc/internal/serverinterceptors"

type (
	ServerMiddlewaresConf struct {
		Trace      bool     `json:",default=true"`
		Recover    bool     `json:",default=true"`
		Stat       bool     `json:",default=true"`
		StatConf   StatConf `json:",optional"`
		Prometheus bool     `json:",default=true"`
		Breaker    bool     `json:",default=true"`
	}

	ClientMiddlewaresConf struct {
		Trace      bool `json:",default=true"`
		Duration   bool `json:",default=true"`
		Prometheus bool `json:",default=true"`
		Breaker    bool `json:",default=true"`
		Timeout    bool `json:",default=true"`
	}

	MethodTimeoutConf = serverinterceptors.MethodTimeoutConf

	StatConf = serverinterceptors.StatConf
)
