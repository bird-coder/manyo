/*
 * @Author: yujiajie
 * @Date: 2025-01-15 16:11:01
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-03-21 16:59:33
 * @FilePath: /manyo/pkg/discov/etcd/clients.go
 * @Description:
 */
package etcd

import (
	"fmt"
	"strings"

	"github.com/bird-coder/manyo/pkg/discov/etcd/internal"
)

const (
	indexOfId = iota + 1
)

const timeToLive int64 = 10

// TimeToLive is seconds to live in etcd.
var TimeToLive = timeToLive

func extract(etcdKey string, index int) (string, bool) {
	if index < 0 {
		return "", false
	}

	fields := strings.FieldsFunc(etcdKey, func(ch rune) bool {
		return ch == internal.Delimiter
	})
	if index >= len(fields) {
		return "", false
	}

	return fields[index], true
}

func extractId(etcdKey string) (string, bool) {
	return extract(etcdKey, indexOfId)
}

func makeEtcdKey(key string, id int64) string {
	return fmt.Sprintf("%s%c%d", key, internal.Delimiter, id)
}
