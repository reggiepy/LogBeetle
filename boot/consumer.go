package boot

import (
	"fmt"
	"github.com/reggiepy/LogBeetle/global"
	"github.com/reggiepy/LogBeetle/pkg/consumer/manager"
	"github.com/reggiepy/LogBeetle/pkg/consumer/nsq_consumer"
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
	c, err := nsq_consumer.NewNSQLogConsumer(
		nsq_consumer.WithName("test"),
		nsq_consumer.WithLogFileName("test.log"),
		nsq_consumer.WithNSQTopic("test"),
		nsq_consumer.WithNSQAddress(nsqConfig.NSQDAddress),
		nsq_consumer.WithNSQAuthSecret(nsqConfig.AuthSecret),
	)
	if err != nil {
		global.LbLogger.Fatal(fmt.Sprintf("error creating consumer %s: %v", "test", err))
	} else {
		global.LBConsumerManager.Add(c)
		global.LbLogger.Info(fmt.Sprintf("consumer %s added to consumer manager", "test"))
	}

	// 添加其他消费者
	for _, cfg := range consumerConfig.NSQConsumers {
		var topic string
		if cfg.Topic == "" {
			topic = cfg.Name
		}
		c, err := nsq_consumer.NewNSQLogConsumer(
			nsq_consumer.WithName(cfg.Name),
			nsq_consumer.WithLogFileName(cfg.FileName),
			nsq_consumer.WithNSQTopic(topic),
			nsq_consumer.WithNSQAddress(nsqConfig.NSQDAddress),
			nsq_consumer.WithNSQAuthSecret(nsqConfig.AuthSecret),
		)
		if err != nil {
			global.LbLogger.Fatal(fmt.Sprintf("error creating consumer %s: %v", cfg.Topic, err))
		} else {
			global.LBConsumerManager.Add(c)
			global.LbLogger.Info(fmt.Sprintf("consumer %s added to consumer manager", cfg.Name))
		}
	}

	global.LBConsumerManager.Start()
	global.OnExit(func() {
		global.LBConsumerManager.Stop()
		global.LbLogger.Info("NSQ consumer stopped")
	})
}
