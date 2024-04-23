package consumer

import (
	"fmt"
	"github.com/reggiepy/LogBeetle/pkg/config"
	"go.uber.org/zap"
	"io"
	"log"
	"path"

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
			panic(err.(any))
		}
	}
	cfg := nsq.NewConfig()
	if n.AuthSecret != "" {
		if err := cfg.Set("auth_secret", n.AuthSecret); err != nil {
			log.Fatalf("Failed to set auth_secret: %v", err)
		}
	}
	if n.Topic == "" {
		panic(fmt.Errorf("topic is not set when creating consumer").(any))
	}
	n.Consumer, err = nsq.NewConsumer(n.Topic, "consumer", cfg)
	if err != nil {
		panic(fmt.Errorf("failed to create NSQ Consumer: %v", err).(any))
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

type NSQLogConsumer struct {
	Name        string
	LogFileName string
	LogFile     io.WriteCloser
	// Logger    zerolog.Logger
	Logger *zap.Logger

	// Consumer
	Consumer *NsqConsumer
}

func (c *NSQLogConsumer) Close() {
	var err error
	if c.LogFile != nil {
		fmt.Printf("关闭 【%s】 Consumer LogFile\n", c.Name)
		err = c.Logger.Sync()
		if err != nil {
			fmt.Printf("Sync 【%s】  Logger 失败: %v", c.Name, err)
		}
		err = c.LogFile.Close()
		if err != nil {
			fmt.Printf("关闭 【%s】 日志文件失败: %v", c.Name, err)
		}
	}
	if c.Consumer != nil {
		fmt.Printf("关闭 【%s】 Consumer Consumer\n", c.Name)
		c.Consumer.Stop()
	}
}

func NewNSQLogConsumer(name string, logFileName string, consumer *NsqConsumer) *NSQLogConsumer {
	c := &NSQLogConsumer{
		Name:        name,
		LogFileName: logFileName,
		Consumer:    consumer,
	}
	// 创建 lumberjack.Logger 实例用于日志切割
	consumerConfig := config.Instance.ConsumerConfig
	filePath := path.Join(consumerConfig.LogPath, c.LogFileName)
	c.LogFile = NewLJLoggerWriteCloser(filePath)

	//logger := NewZEROLogger(c.LogFile)

	// 创建Logger
	c.Logger = NewZAPLogger(c.LogFile)
	c.Consumer.AddHandler(&MessageHandler{
		Handler: func(message []byte) error {
			c.Logger.Info(string(message))
			return nil
		},
	})
	err := c.Consumer.Connect()
	if err != nil {
		panic((err).(any))
	}
	return c
}
