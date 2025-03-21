/*
 * @Author: yujiajie
 * @Date: 2025-01-15 16:11:34
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-01-23 15:13:34
 * @FilePath: /Go-Base/pkg/discov/etcd/internal/vars.go
 * @Description:
 */
package internal

import "time"

const (
	// Delimiter is a separator that separates the etcd path.
	Delimiter = '/'

	autoSyncInterval   = time.Minute
	coolDownInterval   = time.Second
	dialTimeout        = 5 * time.Second
	requestTimeout     = 3 * time.Second
	endpointsSeparator = ","
)

var (
	// DialTimeout is the dial timeout.
	DialTimeout = dialTimeout
	// RequestTimeout is the request timeout.
	RequestTimeout = requestTimeout
	// NewClient is used to create etcd clients.
	// NewClient = DialClient
)
