/*
 * @Author: yujiajie
 * @Date: 2024-12-25 19:43:53
 * @LastEditors: yujiajie 1037297660@qq.com
 * @LastEditTime: 2026-05-12 11:00:08
 * @FilePath: /manyo/pkg/core/config.go
 * @Description:
 */
package core

import (
	"fmt"

	cfg "github.com/bird-coder/manyo/config"
	"github.com/bird-coder/manyo/constant"
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

// Normalize 负责补齐配置结构本身的默认值，避免后续初始化阶段大量出现 nil 分支。
// 这里仅处理“结构完整性”相关的默认值，不处理依赖外部资源的校验逻辑。
func (app *BaseAppConfig) Normalize() {
	if app.System == nil {
		app.System = &SysConfig{}
	}
	if len(app.System.Environment) == 0 {
		app.System.Environment = constant.Dev.String()
	}
	if len(app.System.Timezone) == 0 {
		app.System.Timezone = defaultTimeZone
	}

	if app.Loggers == nil {
		app.Loggers = make(map[string]*cfg.LoggerConfig)
	}
	if app.Databases == nil {
		app.Databases = make(map[string]*cfg.MysqlConfig)
	}
	if app.Redis == nil {
		app.Redis = make(map[string]*cfg.RedisDailConfig)
	}
	if app.RocketMq == nil {
		app.RocketMq = make(map[string]*cfg.MqConfig)
	}
	if app.Consumers == nil {
		app.Consumers = make(map[string]*cfg.ConsumerConfig)
	}
}

// Validate 负责做配置层的静态校验。
// 这里只检查“配置是否自洽”，不检查 redis/mysql/nacos 这类需要真实连接外部资源的能力。
func (app *BaseAppConfig) Validate() error {
	if app.System == nil {
		return fmt.Errorf("system config is required")
	}
	if len(app.System.AppName) == 0 {
		return fmt.Errorf("system.appName is required")
	}

	for key, conf := range app.Loggers {
		if conf == nil {
			return fmt.Errorf("loggers.%s is nil", key)
		}
	}
	for key, conf := range app.Databases {
		if conf == nil {
			return fmt.Errorf("databases.%s is nil", key)
		}
	}
	for key, conf := range app.Redis {
		if conf == nil {
			return fmt.Errorf("redis.%s is nil", key)
		}
	}
	for key, conf := range app.RocketMq {
		if conf == nil {
			return fmt.Errorf("rocketmq.%s is nil", key)
		}
	}
	for key, conf := range app.Consumers {
		if conf == nil {
			return fmt.Errorf("consumers.%s is nil", key)
		}
	}

	return nil
}

func (app *BaseAppConfig) LoadConfig(configFile string) (err error) {
	viper.SetConfigFile(configFile)
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	if err = viper.Unmarshal(&app); err != nil {
		return
	}
	app.Normalize()
	if err = app.Validate(); err != nil {
		return
	}
	return
}
