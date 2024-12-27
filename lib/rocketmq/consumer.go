/*
 * @Author: yujiajie
 * @Date: 2024-05-16 17:28:12
 * @LastEditors: yujiajie
 * @LastEditTime: 2024-12-27 16:25:47
 * @FilePath: /manyo/lib/rocketmq/consumer.go
 * @Description:
 */
package rocketmq

import (
	"context"
	"errors"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
	"github.com/apache/rocketmq-clients/golang/v5/credentials"
	rpcProtoclErr "github.com/apache/rocketmq-clients/golang/v5/protocol/v2"
	"github.com/bird-coder/manyo/config"
	"github.com/bird-coder/manyo/pkg/logger"
)

var TagSeparator = "||"

var (
	ErrConsumerOption = errors.New("consumer option error")
	ErrProductOption  = errors.New("product option error")
)

type MsgHandler interface {
	Handle(mv MsgView) error
}

type MsgView struct {
	Mv *rmq_client.MessageView
	Consumer
}

type Consumer interface {
	Run() error
	Ack(mv *rmq_client.MessageView) error
	Close()
}

type ConsumerGroup struct {
	GroupName string
	Params    []ConsumerParam
}

type ConsumerParam struct {
	TopicName     string
	TagExpression string
	MsgHandler    MsgHandler
}

type NormalConsumer struct {
	opts *options
	rmq_client.SimpleConsumer

	// 外部参数
	ConsumerGroup

	ctx      context.Context
	stopChan chan struct{}
	once     sync.Once
}

func checkGroupParam(group ConsumerGroup) bool {
	if len(group.GroupName) == 0 || len(group.Params) == 0 {
		return false
	}
	for _, param := range group.Params {
		if len(param.TopicName) == 0 || len(param.TagExpression) == 0 || param.MsgHandler == nil {
			return false
		}
	}
	return true
}

func parseExpressions(group ConsumerGroup) map[string]*rmq_client.FilterExpression {
	tmp := make(map[string][]string)
	for _, param := range group.Params {
		if _, exists := tmp[param.TopicName]; !exists {
			tmp[param.TopicName] = make([]string, 0)
		}
		tmp[param.TopicName] = append(tmp[param.TopicName], param.TagExpression)

	}
	expressions := make(map[string]*rmq_client.FilterExpression, len(tmp))
	for k, sl := range tmp {
		expressions[k] = rmq_client.NewFilterExpression(strings.Join(sl, "||"))
	}
	return expressions
}

func initConsumer(opts *options, group ConsumerGroup) (rmq_client.SimpleConsumer, error) {
	// sdk log输出到文件的跟目录
	if len(opts.clientLogRoot) > 0 {
		_ = os.Setenv(rmq_client.CLIENT_LOG_ROOT, opts.clientLogRoot)
	}
	if len(opts.clientLogFilename) > 0 {
		_ = os.Setenv(rmq_client.CLIENT_LOG_FILENAME, opts.clientLogFilename)
	}
	// log to console
	if opts.LogToStdout {
		_ = os.Setenv("mq.consoleAppender.enabled", "true")
	}
	rmq_client.ResetLogger()

	expressions := parseExpressions(group)
	return rmq_client.NewSimpleConsumer(&rmq_client.Config{
		Endpoint:      opts.Endpoint,
		ConsumerGroup: group.GroupName,
		Credentials: &credentials.SessionCredentials{
			AccessKey:    opts.AccessKey,
			AccessSecret: opts.AccessSecret,
		},
	},
		rmq_client.WithAwaitDuration(opts.awaitDuration),
		rmq_client.WithSubscriptionExpressions(expressions),
	)
}

func NewNormalConsumer(ctx context.Context, mqConfig *config.MqConfig, group ConsumerGroup, options ...OptionFunc) (Consumer, error) {
	if !checkGroupParam(group) {
		return nil, ErrConsumerOption
	}
	opts := defaultOptions()
	for _, optionFuc := range options {
		optionFuc(opts)
	}
	opts.MqConfig = mqConfig

	simpleConsumer, err := initConsumer(opts, group)
	if err != nil {
		return nil, err
	}
	mqConsumer := &NormalConsumer{
		opts:           opts,
		ConsumerGroup:  group,
		ctx:            ctx,
		stopChan:       make(chan struct{}),
		SimpleConsumer: simpleConsumer,
	}
	return mqConsumer, nil
}

func (mqConsumer *NormalConsumer) Run() error {
	err := mqConsumer.SimpleConsumer.Start()
	if err != nil {
		return err
	}

	for {
		select {
		case <-mqConsumer.stopChan:
			logger.Info("消费者退出。")
			if err := mqConsumer.SimpleConsumer.GracefulStop(); err != nil {
				mqConsumer.loggerError("mq consumer stop err", err)
			}
			return nil
		default:
			mvs, err := mqConsumer.SimpleConsumer.Receive(mqConsumer.ctx, mqConsumer.opts.maxMessageNum, mqConsumer.opts.invisibleDuration)
			// 只记录非 没有消息的错误
			if err != nil {
				if rpcErr, isRpcError := err.(*rmq_client.ErrRpcStatus); isRpcError {
					if rpcErr.GetCode() != int32(rpcProtoclErr.Code_MESSAGE_NOT_FOUND) {
						mqConsumer.loggerError("mq consumer Receive err", err)
					}
				} else {
					// 这里应该走不到
					mqConsumer.loggerError("mq consumer Receive err2", err)
				}
			}
			// 处理消息，并 ack message
			for _, mv := range mvs {
				tmv := mv
				handler := mqConsumer.getHandler(tmv)
				if handler == nil {
					mqConsumer.loggerError("mq consumer handle not set", err)
					continue
				}
				msgView := MsgView{
					Mv:       tmv,
					Consumer: mqConsumer,
				}
				err = handler.Handle(msgView)
				if err != nil {
					mqConsumer.loggerError("mq consumer handle err", err)
				}
			}
		}
	}
}

func (mqConsumer *NormalConsumer) getHandler(mv *rmq_client.MessageView) MsgHandler {
	for _, param := range mqConsumer.Params {
		tags := strings.Split(strings.Trim(param.TagExpression, TagSeparator), TagSeparator)
		if param.TopicName == mv.GetTopic() && slices.Contains[[]string, string](tags, *mv.GetTag()) {
			return param.MsgHandler
		}
	}
	return nil
}

func (mqConsumer *NormalConsumer) Ack(mv *rmq_client.MessageView) error {
	return mqConsumer.SimpleConsumer.Ack(mqConsumer.ctx, mv)
}

func (mqConsumer *NormalConsumer) Close() {
	mqConsumer.once.Do(func() {
		close(mqConsumer.stopChan)
	})
}

// 记录错误
func (mqConsumer *NormalConsumer) loggerError(msg string, err error) {
	logger.Error("msg: %s, time: %v, err: %v, group: %v", msg, time.Now(), err, mqConsumer.ConsumerGroup)
}
