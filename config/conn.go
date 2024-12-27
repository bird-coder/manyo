/*
 * @Author: yujiajie
 * @Date: 2024-05-31 15:30:47
 * @LastEditors: yujiajie
 * @LastEditTime: 2024-05-31 15:31:03
 * @FilePath: /manyo/config/conn.go
 * @Description:
 */
package config

type ConnConfig struct {
	Addr         string
	ReadTimeout  int
	WriteTimeout int
	MaxMsgBytes  int
}
