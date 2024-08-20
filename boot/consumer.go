package boot

import (
	"fmt"
	"github.com/reggiepy/LogBeetle/global"
	"github.com/reggiepy/LogBeetle/pkg/consumer"
	"github.com/reggiepy/LogBeetle/pkg/consumer/manager"
)

func Consumer() {
	// 获取配置
	global.LBConsumerManager = &manager.Manager{}
	nsqConfig := global.LbConfig.NSQConfig
	consumerConfig := global.LbConfig.ConsumerConfig

	if len(consumerConfig.NSQConsumers) == 0 {
		global.LbLogger.Fatal(fmt.Sprintf("consumer config is empty"))
	}

	// 创建并添加主消费者
	c, err := consumer.NewNSQLogConsumer(
		consumer.WithName("test"),
		consumer.WithLogFileName("test.log"),
		consumer.WithNSQTopic("test"),
		consumer.WithNSQAddress(nsqConfig.NSQDAddress),
		consumer.WithNSQAuthSecret(nsqConfig.AuthSecret),
	)
	if err != nil {
		global.LbLogger.Fatal(fmt.Sprintf("error creating consumer %s: %v", "test", err))
	} else {
		global.LBConsumerManager.Add(c)
		global.LbLogger.Info(fmt.Sprintf("consumer %s added to consumer manager", "test"))
	}

	// 添加其他消费者
	for _, cfg := range consumerConfig.NSQConsumers {
		c, err := consumer.NewNSQLogConsumer(
			consumer.WithName(cfg.Name),
			consumer.WithLogFileName(cfg.FileName),
			consumer.WithNSQTopic(cfg.Topic),
			consumer.WithNSQAddress(nsqConfig.NSQDAddress),
			consumer.WithNSQAuthSecret(nsqConfig.AuthSecret),
		)
		if err != nil {
			global.LbLogger.Fatal(fmt.Sprintf("error creating consumer %s: %v", cfg.Topic, err))
		} else {
			global.LBConsumerManager.Add(c)
			global.LbLogger.Info(fmt.Sprintf("consumer %s added to consumer manager", cfg.Name))
		}
	}

	global.LBConsumerManager.Start()
}
