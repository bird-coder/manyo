/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-10-03 14:13:47
 * @LastEditTime: 2023-10-03 14:15:30
 * @LastEditors: yuanshisan
 */
package config

type RedisPoolConfig struct {
	Idle        int
	Active      int
	IdleTimeout int64
	Wait        bool
}

type RedisDailConfig struct {
	DialTimeout  int64
	ReadTimeout  int64
	WriteTimeout int64
	Protocol     string
	Addr         string
	Db           int
	Password     string
}
