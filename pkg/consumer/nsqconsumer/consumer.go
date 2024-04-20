package nsqconsumer

import (
	"fmt"
	"log"

	"github.com/nsqio/go-nsq"
)

type NsqConsumer struct {
	Topic      string
	Address    string
	AuthSecret string
	Handlers   []*MessageHandler

	Consumer *nsq.Consumer
}

type Options func(n *NsqConsumer) error

// WithAuthSecret 设置NSQ消费者的认证秘钥
func WithAuthSecret(authSecret string) Options {
	return func(n *NsqConsumer) error {
		n.AuthSecret = authSecret
		return nil
	}
}

// WithHandlers 设置NSQ消费者的消息处理器
func WithHandlers(handlers []*MessageHandler) Options {
	return func(n *NsqConsumer) error {
		n.Handlers = append(n.Handlers, handlers...)
		return nil
	}
}

func NewNsqConsumer(topic string, address string, opts ...Options) *NsqConsumer {
	var err error
	n := &NsqConsumer{
		Topic:   topic,
		Address: address,
	}
	for _, opt := range opts {
		err := opt(n)
		if err != nil {
			panic(err)
		}
	}
	cfg := nsq.NewConfig()
	if n.AuthSecret != "" {
		if err := cfg.Set("auth_secret", n.AuthSecret); err != nil {
			log.Fatalf("Failed to set auth_secret: %v", err)
		}
	}
	if n.Topic == "" {
		panic("Topic is not set when creating consumer")
	}
	n.Consumer, err = nsq.NewConsumer(n.Topic, "consumer", cfg)
	if err != nil {
		log.Fatalf("Failed to create NSQ Consumer: %v", err)
	}

	for _, h := range n.Handlers {
		n.AddHandler(h)
	}

	return n
}

func (n *NsqConsumer) AddHandler(handle *MessageHandler) {
	n.Consumer.AddHandler(handle)
}

func (n *NsqConsumer) Connect() error {
	err := n.Consumer.ConnectToNSQD(n.Address)
	if err != nil {
		return fmt.Errorf("连接NSQ失败:%v", err)
	}
	return nil
}

func (n *NsqConsumer) Stop() {
	n.Consumer.Stop()
}
