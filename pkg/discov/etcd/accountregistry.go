/*
 * @Author: yujiajie
 * @Date: 2025-01-23 17:30:50
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-03-21 16:59:28
 * @FilePath: /manyo/pkg/discov/etcd/accountregistry.go
 * @Description:
 */
package etcd

import "github.com/bird-coder/manyo/pkg/discov/etcd/internal"

func RegisterAccount(endpoints []string, user, pass string) {
	internal.AddAccount(endpoints, user, pass)
}

func RegisterTLS(endpoints []string, certFile, certKeyFile, caFile string, insecureSkipVerify bool) error {
	return internal.AddTLS(endpoints, certFile, certKeyFile, caFile, insecureSkipVerify)
}
