/*
 * @Author: yujiajie
 * @Date: 2024-12-26 09:46:43
 * @LastEditors: yujiajie
 * @LastEditTime: 2024-12-26 09:48:24
 * @FilePath: /Go-Base/config/server.go
 * @Description:
 */
package config

type ServerConfig struct {
	Addr     string `mapstructure:"addr"`
	Timeout  int    `mapstructure:"timeout"`
	MaxBytes int    `mapstructure:"maxBytes"`
}
