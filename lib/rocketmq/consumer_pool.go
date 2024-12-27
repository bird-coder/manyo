/*
 * @Author: yujiajie
 * @Date: 2024-05-16 17:42:13
 * @LastEditors: yujiajie
 * @LastEditTime: 2024-12-27 16:25:33
 * @FilePath: /manyo/lib/rocketmq/consumer_pool.go
 * @Description: 支持并发的消费者
 */
package rocketmq

import (
	"context"
	"slices"
	"strings"
	"sync"
	"time"

	rmq_client "github.com/apache/rocketmq-clients/golang/v5"
	rpcProtoclErr "github.com/apache/rocketmq-clients/golang/v5/protocol/v2"
	"github.com/bird-coder/manyo/config"
	"github.com/bird-coder/manyo/pkg/logger"
	"github.com/bird-coder/manyo/pkg/pool/gopool"
)

type ConcurrentConsumer struct {
	opts *options
	rmq_client.SimpleConsumer

	// 外部参数
	ConsumerGroup

	ctx      context.Context
	stopChan chan struct{}
	p        *gopool.Pool
	once     sync.Once
	wg       sync.WaitGroup
}

func NewConcurrentConsumer(ctx context.Context, mqConfig *config.MqConfig, group ConsumerGroup, options ...OptionFunc) (Consumer, error) {
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
	mqConsumer := &ConcurrentConsumer{
		opts:           opts,
		ConsumerGroup:  group,
		ctx:            ctx,
		stopChan:       make(chan struct{}),
		p:              gopool.NewPool(opts.goroutineNum),
		SimpleConsumer: simpleConsumer,
	}
	return mqConsumer, nil
}

func (mqConsumer *ConcurrentConsumer) Run() error {
	err := mqConsumer.SimpleConsumer.Start()
	if err != nil {
		return err
	}

	for {
		select {
		case <-mqConsumer.stopChan:
			logger.Info("等待消费者退出...")
			mqConsumer.wg.Wait()
			logger.Info("消费者退出")
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
				mqConsumer.wg.Add(1)
				mqConsumer.p.Run(func() {
					defer func() {
						mqConsumer.wg.Done()
					}()
					handler := mqConsumer.getHandler(tmv)
					if handler == nil {
						mqConsumer.loggerError("mq consumer handle not set", err)
						return
					}
					msgView := MsgView{
						Mv:       tmv,
						Consumer: mqConsumer,
					}
					err = handler.Handle(msgView)
					if err != nil {
						mqConsumer.loggerError("mq consumer handle err", err)
					}
				})
			}
		}
	}
}

func (mqConsumer *ConcurrentConsumer) Ack(mv *rmq_client.MessageView) error {
	return mqConsumer.SimpleConsumer.Ack(mqConsumer.ctx, mv)
}

func (mqConsumer *ConcurrentConsumer) getHandler(mv *rmq_client.MessageView) MsgHandler {
	for _, param := range mqConsumer.Params {
		tags := strings.Split(strings.Trim(param.TagExpression, TagSeparator), TagSeparator)
		if param.TopicName == mv.GetTopic() && slices.Contains[[]string, string](tags, *mv.GetTag()) {
			return param.MsgHandler
		}
	}
	return nil
}

func (mqConsumer *ConcurrentConsumer) Close() {
	mqConsumer.once.Do(func() {
		close(mqConsumer.stopChan)
	})
}

// 记录错误
func (mqConsumer *ConcurrentConsumer) loggerError(msg string, err error) {
	logger.Error("msg: %s, time: %v, err: %v, group: %v", msg, time.Now(), err, mqConsumer.ConsumerGroup)
}
