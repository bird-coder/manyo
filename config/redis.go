/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-10-03 14:13:47
 * @LastEditTime: 2025-06-04 17:20:08
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
	DialTimeout  int64  `mapstructure:"dialTimeout"`
	ReadTimeout  int64  `mapstructure:"readTimeout"`
	WriteTimeout int64  `mapstructure:"writeTimeout"`
	Protocol     int    `mapstructure:"protocol"`
	Addr         string `mapstructure:"addr"`
	Db           int    `mapstructure:"db"`
	Password     string `mapstructure:"pwd_password"`
	PoolSize     int    `mapstructure:"poolSize"`
	IdleConns    int    `mapstructure:"idleConns"`
	MaxRetry     int    `mapstructure:"maxRetry"`
}
