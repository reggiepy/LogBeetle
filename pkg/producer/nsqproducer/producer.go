package nsqproducer

import (
	"github.com/nsqio/go-nsq"
	"log"
	"sync"
)

var (
	Producer *nsq.Producer
	once     sync.Once
)

type ProducerConfig struct {
	Address    string
	AuthSecret string
}

func InitProducer(config ProducerConfig) {
	once.Do(func() {
		initProducer(config)
	})
}

func initProducer(config ProducerConfig) {
	cfg := nsq.NewConfig()
	if config.AuthSecret != "" {
		if err := cfg.Set("auth_secret", config.AuthSecret); err != nil {
			log.Fatalf("Failed to set auth_secret: %v", err)
		}
	}
	var err error
	Producer, err = nsq.NewProducer(config.Address, cfg)
	if err != nil {
		log.Fatalf("Failed to create NSQ Producer: %v", err)
	}
}

func StopProducer() {
	if Producer != nil {
		Producer.Stop()
	}
}
