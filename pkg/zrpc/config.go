/*
 * @Author: yujiajie
 * @Date: 2025-01-02 17:11:30
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-03-21 17:01:48
 * @FilePath: /manyo/pkg/zrpc/config.go
 * @Description:
 */
package zrpc

import (
	"time"

	"github.com/bird-coder/manyo/pkg/zrpc/internal"
)

type (
	ServerMiddlewaresConf = internal.ServerMiddlewaresConf

	ClientMiddlewaresConf = internal.ClientMiddlewaresConf

	MethodTimeoutConf = internal.MethodTimeoutConf

	StatConf = internal.StatConf

	RpcClientConf struct {
		Target        string        `json:",optional"`
		NonBlock      bool          `json:",optional"`
		Timeout       int64         `json:",default=2000"`
		KeepaliveTime time.Duration `json:",optional"`
		Middlewares   ClientMiddlewaresConf
	}

	RpcServerConf struct {
		Name           string
		Addr           string
		Timeout        int64 `json:",default=2000"`
		Middlewares    ServerMiddlewaresConf
		MethodTimeouts []MethodTimeoutConf `json:",optional"`
	}
)
