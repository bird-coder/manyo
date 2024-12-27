/*
 * @Author: yujiajie
 * @Date: 2024-12-25 19:43:53
 * @LastEditors: yujiajie
 * @LastEditTime: 2024-12-27 16:17:16
 * @FilePath: /manyo/pkg/core/config.go
 * @Description:
 */
package core

import (
	cfg "github.com/bird-coder/manyo/config"
)

type BaseAppConfig struct {
	Name        string                          `yaml:"name"`
	Loggers     map[string]*cfg.LoggerConfig    `mapstructure:"loggers"`
	Databases   map[string]*cfg.MysqlConfig     `yaml:"databases" mapstructure:"databases"`
	Redis       map[string]*cfg.RedisDailConfig `yaml:"redis"`
	Locker      *cfg.RedisDailConfig            `yaml:"locker"`
	RocketMq    map[string]*cfg.MqConfig        `mapstructure:"rocketmq"`
	Consumers   map[string]*cfg.ConsumerConfig  `mapstructure:"consumers"`
	Environment string                          `mapstructure:"environment"`
}
