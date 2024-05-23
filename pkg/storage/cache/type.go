/*
 * @Author: yujiajie
 * @Date: 2024-05-14 09:55:53
 * @LastEditors: yujiajie
 * @LastEditTime: 2024-05-23 17:56:14
 * @FilePath: /manyo/pkg/storage/cache/type.go
 * @Description:
 */
package cache

import "time"

type AdapterCache interface {
	String() string
	Get(key string) (string, error)
	Set(key string, val interface{}, expire int) error
	Del(key string) error
	HGet(hk, key string) (string, error)
	HDel(hk, key string) error
	Increase(key string) error
	Decrease(key string) error
	Expire(key string, dur time.Duration) error
}
