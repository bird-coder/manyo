/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-10-03 14:13:47
 * @LastEditTime: 2024-03-16 22:47:24
 * @LastEditors: yujiajie
 */
package config

type CacheConfig struct {
	Redis *RedisDailConfig
}

type LockConfig struct {
	Redis *RedisDailConfig
}

type RedisDailConfig struct {
	DialTimeout  int64
	ReadTimeout  int64
	WriteTimeout int64
	Protocol     int
	Addr         string
	Db           int
	Password     string
	PoolSize     int
	IdleConns    int
	MaxRetry     int
}
