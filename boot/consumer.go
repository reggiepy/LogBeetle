package boot

import (
	"fmt"
	"github.com/reggiepy/LogBeetle/global"
	"github.com/reggiepy/LogBeetle/pkg/consumer/manager"
	"github.com/reggiepy/LogBeetle/pkg/consumer/nsq_consumer"
	"github.com/reggiepy/goutils/signailUtils"
)

func Consumer() *manager.Manager {
	// 获取配置
	consumerManager, err := manager.NewManager(manager.WithLogger(global.LbLogger))
	if err != nil {
		panic(err)
	}
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
		err = consumerManager.Add(c)
		if err != nil {
			global.LbLogger.Fatal(fmt.Sprintf("add consumer error: %v", err))
		} else {
			global.LbLogger.Info(fmt.Sprintf("consumer %s added to consumer manager", "test"))
		}
	}

	// 添加其他消费者
	for _, cfg := range consumerConfig.NSQConsumers {
		var topic = cfg.Topic
		if topic == "" {
			topic = cfg.Name
		}
		if topic == "test" {
			global.LbLogger.Warn("consumer topic can't be 'test'")
			continue
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
			err = consumerManager.Add(c)
			if err != nil {
				global.LbLogger.Fatal(fmt.Sprintf("add consumer error: %v", err))
			} else {
				global.LbLogger.Info(fmt.Sprintf("consumer %s added to consumer manager", "test"))
			}
		}
	}

	consumerManager.Start()
	signailUtils.OnExit(func() {
		consumerManager.Stop()
		global.LbLogger.Info("NSQ consumer stopped")
	})
	return consumerManager
}
