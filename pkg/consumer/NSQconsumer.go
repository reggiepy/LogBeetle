package consumer

import (
	"fmt"
	"github.com/reggiepy/LogBeetle/global"
	"io"
	"path"

	"github.com/nsqio/go-nsq"
	"github.com/reggiepy/LogBeetle/util/struct_utils"
	"go.uber.org/zap"
)

var Topics []string

type NSQLogConsumer struct {
	Name        string
	LogFileName string
	LogFile     io.WriteCloser
	// Logger    zerolog.Logger
	Logger *zap.Logger

	// Consumer
	NSQTopic      string
	NSQAddress    string
	NSQAuthSecret string
	NSQConsumer   *nsq.Consumer
}

func (c *NSQLogConsumer) Close() error {
	var err error
	if c.LogFile != nil {
		fmt.Printf("关闭 【%s】 Consumer LogFile\n", c.Name)
		err = c.Logger.Sync()
		if err != nil {
			return fmt.Errorf("Sync 【%s】  Logger 失败: %v", c.Name, err)
		}
		err = c.LogFile.Close()
		if err != nil {
			return fmt.Errorf("关闭 【%s】 日志文件失败: %v", c.Name, err)
		}
	}
	if c.NSQConsumer != nil {
		fmt.Printf("关闭 【%s】 Consumer NSQConsumer\n", c.Name)
		c.NSQConsumer.Stop()
	}
	return nil
}

func (c *NSQLogConsumer) RegisterTopic(topic ...string) {
	Topics = append(Topics, topic...)
}

type Options func(n *NSQLogConsumer) error

// WithAuthSecret 设置NSQ消费者的认证秘钥
func WithNSQAuthSecret(authSecret string) Options {
	return func(n *NSQLogConsumer) error {
		n.NSQAuthSecret = authSecret
		return nil
	}
}

// WithName 设置NSQ消费者的名称
func WithName(name string) Options {
	return func(n *NSQLogConsumer) error {
		n.Name = name
		return nil
	}
}

// WithLogFileName 设置NSQ消费者的日志文件名
func WithLogFileName(logFileName string) Options {
	return func(n *NSQLogConsumer) error {
		n.LogFileName = logFileName
		return nil
	}
}

// WithLogFile 设置NSQ消费者的日志文件
func WithLogFile(logFile io.WriteCloser) Options {
	return func(n *NSQLogConsumer) error {
		n.LogFile = logFile
		return nil
	}
}

// WithLogger 设置NSQ消费者的日志记录器
func WithLogger(logger *zap.Logger) Options {
	return func(n *NSQLogConsumer) error {
		n.Logger = logger
		return nil
	}
}

// WithNSQTopic 设置NSQ消费者的主题
func WithNSQTopic(topic string) Options {
	return func(n *NSQLogConsumer) error {
		n.NSQTopic = topic
		return nil
	}
}

// WithNSQAddress 设置NSQ消费者的地址
func WithNSQAddress(address string) Options {
	return func(n *NSQLogConsumer) error {
		n.NSQAddress = address
		return nil
	}
}

func NewNSQLogConsumer(opts ...Options) (*NSQLogConsumer, error) {
	c := &NSQLogConsumer{}
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			panic(err.(any))
		}
	}
	status, err := struct_utils.IsEmptyStringField(c, "NSQTopic", "NSQAddress", "LogFileName", "Name")
	if err != nil {
		return nil, err
	}
	if status {
		return nil, fmt.Errorf("%v", err)
	}
	cfg := nsq.NewConfig()
	if c.NSQAuthSecret != "" {
		if err := cfg.Set("auth_secret", c.NSQAuthSecret); err != nil {
			return nil, fmt.Errorf("failed to set auth_secret: %v", err)
		}
	}

	consumer, err := nsq.NewConsumer(c.NSQTopic, "consumer", cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create NSQ Consumer: %v", err)
	}
	c.NSQConsumer = consumer
	// 创建 lumberjack.Logger 实例用于日志切割
	consumerConfig := global.LbConfig.ConsumerConfig
	filePath := path.Join(consumerConfig.LogPath, c.LogFileName)
	c.LogFile = NewLJLoggerWriteCloser(filePath)

	// logger := NewZEROLogger(c.LogFile)

	// 创建Logger
	c.Logger = NewZAPLogger(c.LogFile)
	c.NSQConsumer.AddHandler(&MessageHandler{
		Handler: func(message []byte) error {
			c.Logger.Info(string(message))
			return nil
		},
	})
	err = c.NSQConsumer.ConnectToNSQD(c.NSQAddress)
	if err != nil {
		return nil, fmt.Errorf("连接NSQ失败:%v", err)
	}
	c.RegisterTopic(c.NSQTopic)
	return c, nil
}
