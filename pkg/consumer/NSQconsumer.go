package consumer

import (
	"fmt"
	"io"
	"path"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/reggiepy/LogBeetle/global"

	"github.com/nsqio/go-nsq"
	"github.com/reggiepy/LogBeetle/goutils/structUtils"
	"go.uber.org/zap"
)

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
		err = c.Logger.Sync()
		if err != nil {
			return fmt.Errorf("sync 【%s】  Logger faild: %v", c.Name, err)
		}
		global.LbLogger.Info(fmt.Sprintf("sync consumer 【%s】 logger", c.Name))
		err = c.LogFile.Close()
		if err != nil {
			return fmt.Errorf("closing 【%s】 logfile faild: %v", c.Name, err)
		}
		global.LbLogger.Info(fmt.Sprintf("closing consumer 【%s】 logger file", c.Name))
	}
	if c.NSQConsumer != nil {
		c.NSQConsumer.Stop()
		global.LbLogger.Info(fmt.Sprintf("closing consumer 【%s】 nsq", c.Name))
	}
	return nil
}

// initializeNSQConsumer initializes the NSQ consumer.
func (c *NSQLogConsumer) initializeNSQConsumer() error {
	if c.NSQTopic == "" || c.NSQAddress == "" || c.LogFileName == "" || c.Name == "" {
		return fmt.Errorf("missing required fields")
	}

	cfg := nsq.NewConfig()
	if c.NSQAuthSecret != "" {
		if err := cfg.Set("auth_secret", c.NSQAuthSecret); err != nil {
			return fmt.Errorf("failed to set auth_secret: %v", err)
		}
	}

	consumer, err := nsq.NewConsumer(c.NSQTopic, "consumer", cfg)
	if err != nil {
		return fmt.Errorf("failed to create NSQ consumer: %v", err)
	}
	c.NSQConsumer = consumer

	return nil
}

// initializeLogger initializes the logger and log file.
func (c *NSQLogConsumer) initializeLogger() error {
	consumerConfig := global.LbConfig.ConsumerConfig
	filePath := path.Join(consumerConfig.LogPath, c.LogFileName)

	lj := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    1, // MB
		MaxBackups: 5,
		MaxAge:     30, // days
		Compress:   false,
	}
	logger := NewZAPLogger(lj)
	c.Logger = logger
	c.LogFile = lj

	c.NSQConsumer.AddHandler(&MessageHandler{
		Handler: func(message []byte) error {
			c.Logger.Info(string(message))
			return nil
		},
	})

	return nil
}

// connectToNSQ connects the consumer to the NSQ.
func (c *NSQLogConsumer) connectToNSQ() error {
	if err := c.NSQConsumer.ConnectToNSQD(c.NSQAddress); err != nil {
		return fmt.Errorf("failed to connect to NSQ: %v", err)
	}
	c.RegisterTopic(c.NSQTopic)
	return nil
}

func (c *NSQLogConsumer) RegisterTopic(topic ...string) {
	global.LbRegisterTopic = append(global.LbRegisterTopic, topic...)
}

func NewNSQLogConsumer(opts ...Options) (*NSQLogConsumer, error) {
	c := &NSQLogConsumer{}
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			panic(err.(any))
		}
	}
	status, err := structUtils.IsEmptyStringField(c, "NSQTopic", "NSQAddress", "LogFileName", "Name")
	if err != nil {
		return nil, err
	}
	if status {
		return nil, fmt.Errorf("%v", err)
	}

	// Validate and initialize NSQ consumer
	if err := c.initializeNSQConsumer(); err != nil {
		return nil, err
	}

	// Initialize logging
	if err := c.initializeLogger(); err != nil {
		return nil, err
	}

	// Register topic and connect to NSQ
	if err := c.connectToNSQ(); err != nil {
		return nil, err
	}

	c.RegisterTopic(c.NSQTopic)
	return c, nil
}
