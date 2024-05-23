/*
 * @Author: yujiajie
 * @Date: 2024-05-14 09:54:54
 * @LastEditors: yujiajie
 * @LastEditTime: 2024-05-23 17:56:42
 * @FilePath: /manyo/pkg/storage/locker/redis.go
 * @Description:
 */
package locker

import (
	"context"
	"time"

	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

func NewRedis(client *redis.Client, options *redis.Options) (*Redis, error) {
	if client == nil {
		client = redis.NewClient(options)
	}
	r := &Redis{
		client: client,
	}
	if err := r.testConnect(); err != nil {
		return nil, err
	}
	return r, nil
}

type Redis struct {
	client *redis.Client
	mutex  *redislock.Client
}

func (r *Redis) String() string {
	return "redis"
}

func (r *Redis) testConnect() error {
	_, err := r.client.Ping(context.TODO()).Result()
	return err
}

func (r *Redis) Lock(key string, ttl int64, options *redislock.Options) (*redislock.Lock, error) {
	if r.mutex == nil {
		r.mutex = redislock.New(r.client)
	}
	return r.mutex.Obtain(context.TODO(), key, time.Duration(ttl)*time.Second, options)
}
