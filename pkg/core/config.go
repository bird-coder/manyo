/*
 * @Author: yujiajie
 * @Date: 2024-12-25 19:43:53
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-06-04 17:58:34
 * @FilePath: /manyo/pkg/core/config.go
 * @Description:
 */
package core

import (
	cfg "github.com/bird-coder/manyo/config"
	"github.com/bird-coder/manyo/pkg/discov/nacos"
)

type BaseAppConfig struct {
	System    *SysConfig                      `mapstructure:"system"`
	Loggers   map[string]*cfg.LoggerConfig    `mapstructure:"loggers"`
	Databases map[string]*cfg.MysqlConfig     `yaml:"databases" mapstructure:"databases"`
	Redis     map[string]*cfg.RedisDailConfig `yaml:"redis"`
	Locker    *cfg.RedisDailConfig            `yaml:"locker"`
	RocketMq  map[string]*cfg.MqConfig        `mapstructure:"rocketmq"`
	Consumers map[string]*cfg.ConsumerConfig  `mapstructure:"consumers"`
	Nacos     *nacos.NacosConf                `yaml:"nacos"`
}

type SysConfig struct {
	AppName     string `mapstructure:"appName"`
	Environment string `mapstructure:"env"`
	Timezone    string `mapstructure:"timezone"`
	ServerId    uint8  `mapstructure:"serverId"`
}
