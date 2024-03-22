package nsqconsumer

import (
	"github.com/nsqio/go-nsq"
	"log"
)

type NsqConsumerConfig struct {
	Topic      string
	Channel    string
	Address    string
	AuthSecret string
	Handler    *MessageHandler
}

func NewNsqConsumer(config NsqConsumerConfig) *nsq.Consumer {
	cfg := nsq.NewConfig()
	if config.AuthSecret != "" {
		if err := cfg.Set("auth_secret", config.AuthSecret); err != nil {
			log.Fatalf("Failed to set auth_secret: %v", err)
		}
	}
	consumer, err := nsq.NewConsumer(config.Topic, config.Channel, cfg)
	if err != nil {
		log.Fatalf("Failed to create NSQ Consumer: %v", err)
	}
	consumer.AddHandler(config.Handler)

	err = consumer.ConnectToNSQD(config.Address)
	if err != nil {
		log.Fatalf("Failed to connect to NSQD: %v", err)
	}
	return consumer
}
