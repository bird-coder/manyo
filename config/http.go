/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-10-03 14:17:11
 * @LastEditTime: 2023-10-03 14:18:10
 * @LastEditors: yuanshisan
 */
package config

type HttpConfig struct {
	Addr           string
	ReadTimeout    int
	WriteTimeout   int
	MaxHeaderBytes int
}
