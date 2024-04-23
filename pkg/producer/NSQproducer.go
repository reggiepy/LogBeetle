package producer

import (
	"log"
	"sync"

	"github.com/nsqio/go-nsq"
)

var (
	NSQProducer *nsq.Producer
	once        sync.Once
)

type NSQProducerConfig struct {
	Address    string
	AuthSecret string
}

func InitNSQProducer(config NSQProducerConfig) {
	once.Do(func() {
		initProducer(config)
	})
}

func initProducer(config NSQProducerConfig) {
	cfg := nsq.NewConfig()
	if config.AuthSecret != "" {
		if err := cfg.Set("auth_secret", config.AuthSecret); err != nil {
			log.Fatalf("Failed to set auth_secret: %v", err)
		}
	}
	var err error
	NSQProducer, err = nsq.NewProducer(config.Address, cfg)
	if err != nil {
		log.Fatalf("Failed to create NSQ Producer: %v", err)
	}
}

func StopProducer() {
	if NSQProducer != nil {
		NSQProducer.Stop()
	}
}
