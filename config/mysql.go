/*
 * @Description:
 * @Author: yuanshisan
 * @Date: 2023-10-03 14:08:06
 * @LastEditTime: 2023-10-03 14:13:18
 * @LastEditors: yuanshisan
 */
package config

type MysqlConfig struct {
	IdleConns    int
	OpenConns    int
	IdleTimeout  int64
	AliveTimeout int64
	Cluster      bool
	Default      *DbConfig
	Sources      []*DbConfig
	Replicas     []*DbConfig
}

type DbConfig struct {
	Protocol string
	Host     string
	Port     string
	User     string
	Pass     string
	Database string
	Charset  string
	Prefix   string
}
