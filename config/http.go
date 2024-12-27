/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-10-03 14:17:11
 * @LastEditTime: 2024-12-27 15:58:28
 * @LastEditors: yujiajie
 */
package config

type HttpConfig struct {
	Addr           string `mapstructure:"addr"`
	ReadTimeout    int    `mapstructure:"readTimeout"`
	WriteTimeout   int    `mapstructure:"writeTimeout"`
	MaxHeaderBytes int    `mapstructure:"maxHeaderBytes"`
}
