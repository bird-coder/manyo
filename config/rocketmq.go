/*
 * @Author: yujiajie
 * @Date: 2024-12-27 15:59:05
 * @LastEditors: yujiajie
 * @LastEditTime: 2024-12-27 15:59:07
 * @FilePath: /manyo/config/rocketmq.go
 * @Description:
 */
package config

type MqConfig struct {
	Endpoint     string `mapstructure:"endpoint"`
	AccessKey    string `mapstructure:"pwd_access_key"`
	AccessSecret string `mapstructure:"pwd_secret_key"`
	LogToStdout  bool   `mapstructure:"log_to_stdout"`
}

type ConsumerConfig struct {
	GroupName string `mapstructure:"groupName"`
	Topic     string `mapstructure:"topic"`
	Tags      string `mapstructure:"tags"`
	LogDir    string `mapstructure:"logDir"`
	PoolSize  int    `mapstructure:"poolSize"`
	Case      int    `mapstructure:"case"`
	MsgNum    int    `mapstructure:"msgNum"`
}
