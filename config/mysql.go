/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-10-03 14:08:06
 * @LastEditTime: 2025-06-04 17:19:23
 * @LastEditors: yujiajie
 */
package config

type MysqlConfig struct {
	Driver       string            `mapstructure:"driver"`
	Dsn          string            `mapstructure:"dsn"`
	Password     string            `mapstructure:"pwd_password"`
	IdleConns    int               `mapstructure:"idleConns"`
	OpenConns    int               `mapstructure:"openConns"`
	IdleTimeout  int64             `mapstructure:"idleTimeout"`
	AliveTimeout int64             `mapstructure:"aliveTimeout"`
	Cluster      bool              `mapstructure:"cluster"`
	Sources      []MysqlBaseConfig `mapstructure:"sources"`
	Replicas     []MysqlBaseConfig `mapstructure:"replicas"`
}

type MysqlBaseConfig struct {
	Dsn      string `mapstructure:"dsn"`
	Password string `mapstructure:"pwd_password"`
}
