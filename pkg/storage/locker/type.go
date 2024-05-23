/*
 * @Author: yujiajie
 * @Date: 2024-05-14 09:55:22
 * @LastEditors: yujiajie
 * @LastEditTime: 2024-05-23 17:56:50
 * @FilePath: /manyo/pkg/storage/locker/type.go
 * @Description:
 */
package locker

import "github.com/bsm/redislock"

type AdapterLocker interface {
	String() string
	Lock(key string, ttl int64, options *redislock.Options) (*redislock.Lock, error)
}
