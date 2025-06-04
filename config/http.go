/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-10-03 14:17:11
 * @LastEditTime: 2025-06-04 17:18:56
 * @LastEditors: yujiajie
 */
package config

type HttpConfig struct {
	Addr           string `mapstructure:"addr"`
	ReadTimeout    int    `mapstructure:"readTimeout"`
	WriteTimeout   int    `mapstructure:"writeTimeout"`
	MaxHeaderBytes int    `mapstructure:"maxHeaderBytes"`
	CertFile       string `mapstructure:"certFile"`
	KeyFile        string `mapstructure:"keyFile"`
}
