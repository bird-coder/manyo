/*
 * @Author: yujiajie
 * @Date: 2025-01-23 14:47:38
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-01-23 14:48:33
 * @FilePath: /Go-Base/pkg/discov/etcd/internal/updatelistener.go
 * @Description:
 */
package internal

type (
	KV struct {
		Key string
		Val string
	}

	UpdateListener interface {
		OnAdd(kv KV)
		OnDelete(kv KV)
	}
)
