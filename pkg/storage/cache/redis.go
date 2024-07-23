/*
 * @Author: yujiajie
 * @Date: 2024-05-14 09:55:38
 * @LastEditors: yujiajie
 * @LastEditTime: 2024-07-23 15:16:09
 * @FilePath: /manyo/pkg/storage/cache/redis.go
 * @Description:
 */
package cache

import (
	"context"
	"time"

	"github.com/bird-coder/manyo/config"
	"github.com/redis/go-redis/v9"
)

func NewRedis(client *redis.Client, cfg *config.RedisDailConfig) (*Redis, error) {
	if client == nil {
		options := &redis.Options{
			Addr:         cfg.Addr,
			Password:     cfg.Password,
			DB:           cfg.Db,
			Protocol:     cfg.Protocol,
			DialTimeout:  time.Duration(cfg.DialTimeout) * time.Second,
			ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
			PoolSize:     cfg.PoolSize,
			MinIdleConns: cfg.IdleConns,
		}
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
}

func (r *Redis) String() string {
	return "redis"
}

func (r *Redis) testConnect() error {
	_, err := r.client.Ping(context.TODO()).Result()
	return err
}

func (r *Redis) Get(key string) (string, error) {
	return r.client.Get(context.TODO(), key).Result()
}

func (r *Redis) Set(key string, val interface{}, expire int) error {
	return r.client.Set(context.TODO(), key, val, time.Duration(expire)*time.Second).Err()
}

func (r *Redis) Del(key string) error {
	return r.client.Del(context.TODO(), key).Err()
}

func (r *Redis) HGet(hk, key string) (string, error) {
	return r.client.HGet(context.TODO(), hk, key).Result()
}

func (r *Redis) HDel(hk, key string) error {
	return r.client.HDel(context.TODO(), hk, key).Err()
}

func (r *Redis) Increase(key string) error {
	return r.client.Incr(context.TODO(), key).Err()
}

func (r *Redis) Decrease(key string) error {
	return r.client.Decr(context.TODO(), key).Err()
}

func (r *Redis) Expire(key string, dur time.Duration) error {
	return r.client.Expire(context.TODO(), key, dur).Err()
}

func (r *Redis) GetClient() *redis.Client {
	return r.client
}
