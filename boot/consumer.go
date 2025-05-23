package boot

import (
	"fmt"
	"github.com/reggiepy/LogBeetle/global"
	"github.com/reggiepy/LogBeetle/pkg/consumer/manager"
	"github.com/reggiepy/LogBeetle/pkg/consumer/nsq_consumer"
	"github.com/reggiepy/goutils/signailUtils"
	"go.uber.org/zap"
)

func Consumer() *manager.Manager {
	// 获取配置
	logger := zap.L().Named("consumer")
	consumerManager, err := manager.NewManager(manager.WithLogger(logger))
	if err != nil {
		panic(err)
	}
	nsqConfig := global.LbConfig.NSQConfig
	consumerConfig := global.LbConfig.ConsumerConfig

	if len(consumerConfig.NSQConsumers) == 0 {
		zap.L().Fatal(fmt.Sprintf("consumer config is empty"))
	}
	options := []nsq_consumer.Options{
		nsq_consumer.WithNSQAuthSecret(nsqConfig.AuthSecret),
		nsq_consumer.WithLogger(logger),
	}

	// 创建并添加主消费者
	c, err := nsq_consumer.NewNSQLogConsumer(
		"test",
		nsqConfig.NSQDAddress,
		options...,
	)
	if err != nil {
		zap.L().Fatal(fmt.Sprintf("error creating consumer %s: %v", "test", err))
	} else {
		err = consumerManager.Add(c)
		if err != nil {
			zap.L().Fatal(fmt.Sprintf("add consumer error: %v", err))
		} else {
			zap.L().Info(fmt.Sprintf("consumer %s added to consumer manager", "test"))
		}
	}

	// 添加其他消费者
	for _, cfg := range consumerConfig.NSQConsumers {
		if cfg.Topic == "test" {
			zap.L().Warn("consumer topic can't be 'test'")
			continue
		}
		options := []nsq_consumer.Options{
			nsq_consumer.WithNSQAuthSecret(nsqConfig.AuthSecret),
			nsq_consumer.WithLogger(logger),
		}
		if cfg.Topic != "" {
			options = append(options, nsq_consumer.WithNSQTopic(cfg.Topic))
		}
		if cfg.FileName != "" {
			options = append(options, nsq_consumer.WithLogFileName(cfg.FileName))
		}
		c, err := nsq_consumer.NewNSQLogConsumer(
			cfg.Name,
			nsqConfig.NSQDAddress,
			options...,
		)
		if err != nil {
			zap.L().Fatal(fmt.Sprintf("error creating consumer %s: %v", cfg.Topic, err))
		} else {
			err = consumerManager.Add(c)
			if err != nil {
				zap.L().Fatal(fmt.Sprintf("add consumer error: %v", err))
			} else {
				zap.L().Info(fmt.Sprintf("consumer %s added to consumer manager", "test"))
			}
		}
	}

	consumerManager.Start()
	signailUtils.OnExit(func() {
		consumerManager.Stop()
		zap.L().Info("NSQ consumer stopped")
	})
	return consumerManager
}
