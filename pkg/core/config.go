/*
 * @Author: yujiajie
 * @Date: 2024-12-25 19:43:53
 * @LastEditors: yujiajie
 * @LastEditTime: 2025-09-03 18:11:18
 * @FilePath: /manyo/pkg/core/config.go
 * @Description:
 */
package core

import (
	cfg "github.com/bird-coder/manyo/config"
	"github.com/bird-coder/manyo/pkg/discov/nacos"

	"github.com/spf13/viper"
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

func (app *BaseAppConfig) LoadConfig(configFile string) (err error) {
	viper.SetConfigFile(configFile)
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	if err = viper.Unmarshal(&app); err != nil {
		return
	}
	Kernal.SetSysInfo(app.System)
	Kernal.SetConfig(CONFIG_KEY_LOGGER, app.Loggers)
	Kernal.SetConfig(CONFIG_KEY_DATABASE, app.Databases)
	Kernal.SetConfig(CONFIG_KEY_REDIS, app.Redis)
	Kernal.SetConfig(CONFIG_KEY_LOCKER, app.Locker)
	Kernal.SetConfig(CONFIG_KEY_ROCKET, app.RocketMq)
	Kernal.SetConfig(CONFIG_KEY_CONSUMER, app.Consumers)
	Kernal.SetConfig(CONFIG_KEY_NACOS, app.Nacos)
	return
}
