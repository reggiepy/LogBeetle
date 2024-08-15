package producer

import (
	"fmt"
	"log"
	"sync"

	"github.com/nsqio/go-nsq"
	"github.com/reggiepy/LogBeetle/goutils/structUtils"
)

var (
	Instance LogBeetleProducer
	once     sync.Once
)

func InitInstance(instance LogBeetleProducer) {
	once.Do(func() {
		Instance = instance
	})
}

type NSQProducer struct {
	NSQAddress    string
	NSQAuthSecret string
	Producer      *nsq.Producer
}

func (p *NSQProducer) Close() error {
	fmt.Printf("关闭 NSQProducer\n")
	p.Producer.Stop()
	return nil
}

func (p *NSQProducer) Publish(topic string, body []byte) error {
	return p.Producer.Publish(topic, body)
}

type Options func(n *NSQProducer) error

// WithAuthSecret 设置NSQ消费者的认证秘钥
func WithNSQAuthSecret(authSecret string) Options {
	return func(n *NSQProducer) error {
		n.NSQAuthSecret = authSecret
		return nil
	}
}

// WithNSQAddress 设置NSQ消费者的地址
func WithNSQAddress(address string) Options {
	return func(n *NSQProducer) error {
		n.NSQAddress = address
		return nil
	}
}

func NewNSQProducer(opts ...Options) (*NSQProducer, error) {
	p := &NSQProducer{}
	for _, opt := range opts {
		err := opt(p)
		if err != nil {
			panic(err.(any))
		}
	}
	status, err := structUtils.IsEmptyStringField(p, "NSQAddress")
	if err != nil {
		return nil, err
	}
	if status {
		return nil, fmt.Errorf("%v", err)
	}
	cfg := nsq.NewConfig()
	if p.NSQAuthSecret != "" {
		if err := cfg.Set("auth_secret", p.NSQAuthSecret); err != nil {
			return nil, fmt.Errorf("failed to set auth_secret: %v", err)
		}
	}
	producer, err := nsq.NewProducer(p.NSQAddress, cfg)
	if err != nil {
		log.Fatalf("Failed to create NSQ Producer: %v", err)
	}

	p.Producer = producer
	return p, nil
}
