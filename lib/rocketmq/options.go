/**
* @Author: maozhongyu
* @Desc: 选项，支持可传， 简化选项， 消费者和生产者共用
* @Date: 2024/4/15
**/
package rocketmq

import (
	"time"

	"github.com/bird-coder/manyo/config"
)

type options struct {
	// 下列项 消费者可能用
	awaitDuration     time.Duration
	maxMessageNum     int32
	invisibleDuration time.Duration
	clientLogFilename string //日志文件
	clientLogRoot     string //日志文件目录
	goroutineNum      int    //最大处理协程数

	// 生产和消费 共用选项参数
	*config.MqConfig
}

const (
	// maximum waiting time for receive func
	defaultAwaitDuration = time.Second * 5
	// maximum number of messages received at one time
	defaultMaxMessageNum int32 = 16
	// invisibleDuration should > 20s
	defaultInvisibleDuration = time.Second * 20
	// 默认的日志文件名
	defaultClientLogFilename = "rocketmq_client_go.log"
	// 默认的日志目录
	defaultClientLogRoot = "./logs/rocketmqlogs"
	//默认最大处理协程数
	defaultGoroutineNum = 2000
)

type OptionFunc func(*options)

// maximum waiting time for receive func
func WithAwaitDuration(t time.Duration) OptionFunc {
	return func(o *options) {
		o.awaitDuration = t
	}
}

// maximum number of messages received at one time
func WithMaxMessageNum(n int32) OptionFunc {
	return func(o *options) {
		o.maxMessageNum = n
	}
}

// invisibleDuration
func WithInvisibleDuration(t time.Duration) OptionFunc {
	return func(o *options) {
		o.invisibleDuration = t
	}
}

// endpoint 支持外面带
func WithEndpoint(e string) OptionFunc {
	return func(o *options) {
		o.Endpoint = e
	}
}

// accessKey
func WithAccessKey(a string) OptionFunc {
	return func(o *options) {
		o.AccessKey = a
	}
}

// accessSecret
func WithAccessSecret(a string) OptionFunc {
	return func(o *options) {
		o.AccessSecret = a
	}
}

// sdk log 是否输出到标准输出
func WithLog2stdout(b bool) OptionFunc {
	return func(o *options) {
		o.LogToStdout = b
	}
}

func WithClientLogFilename(f string) OptionFunc {
	return func(o *options) {
		o.clientLogFilename = f
	}
}

func WithClientLogRoot(f string) OptionFunc {
	return func(o *options) {
		o.clientLogRoot = f
	}
}

func WithGoroutineNum(num int) OptionFunc {
	return func(o *options) {
		o.goroutineNum = num
	}
}

// 默认值
// defaultOptions
//
//	@Description:
//	@return *options , endpoint、access_key、secret_key 默认值从配置读
func defaultOptions() *options {
	opts := &options{
		awaitDuration:     defaultAwaitDuration,
		maxMessageNum:     defaultMaxMessageNum,
		invisibleDuration: defaultInvisibleDuration,
		clientLogFilename: defaultClientLogFilename,
		clientLogRoot:     defaultClientLogRoot,
		goroutineNum:      defaultGoroutineNum,
	}
	return opts
}
