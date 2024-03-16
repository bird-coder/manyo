/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-10-03 14:08:06
 * @LastEditTime: 2024-03-16 22:48:42
 * @LastEditors: yujiajie
 */
package config

type MysqlConfig struct {
	Driver       string   `mapstructure:"driver"`
	IdleConns    int      `mapstructure:"idleConns"`
	OpenConns    int      `mapstructure:"openConns"`
	IdleTimeout  int64    `mapstructure:"idleTimeout"`
	AliveTimeout int64    `mapstructure:"aliveTimeout"`
	Cluster      bool     `mapstructure:"cluster"`
	Default      string   `mapstructure:"default"`
	Sources      []string `mapstructure:"sources"`
	Replicas     []string `mapstructure:"replicas"`
}
