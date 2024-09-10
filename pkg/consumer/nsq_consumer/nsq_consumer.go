package nsq_consumer

import (
	"fmt"
	"io"
	"path"
	"sync"
	"time"

	"go.uber.org/multierr"

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
	Logger      *zap.Logger

	stopOnce  sync.Once
	startOnce sync.Once
	isStarted bool

	// Consumer
	NSQTopic      string
	NSQAddress    string
	NSQAuthSecret string
	NSQConsumer   *nsq.Consumer
	Handler       func(message []byte) error
}

func (c *NSQLogConsumer) Stop() error {
	var err error

	// 确保 Close 只执行一次
	c.stopOnce.Do(func() {
		if c.NSQConsumer != nil {
			c.NSQConsumer.Stop()
			global.LbLogger.Info(fmt.Sprintf("closing consumer 【%s】 nsq", c.Name))
		}

		if c.LogFile != nil {
			// 尝试同步日志
			if syncErr := c.Logger.Sync(); syncErr != nil {
				err = multierr.Append(err, fmt.Errorf("sync 【%s】 Logger failed: %v", c.Name, syncErr))
			} else {
				global.LbLogger.Info(fmt.Sprintf("sync consumer 【%s】 logger", c.Name))
			}

			// 尝试关闭日志文件
			if closeErr := c.LogFile.Close(); closeErr != nil {
				err = multierr.Append(err, fmt.Errorf("closing 【%s】 logfile failed: %v", c.Name, closeErr))
			} else {
				global.LbLogger.Info(fmt.Sprintf("closing consumer 【%s】 logger file", c.Name))
			}
		}
	})
	return err
}

func (c *NSQLogConsumer) Topic() string {
	return c.NSQTopic
}

// Start nsq consumer.
func (c *NSQLogConsumer) Start() error {
	go func() {
		c.startOnce.Do(func() {
			for {
				err := c.NSQConsumer.ConnectToNSQD(c.NSQAddress)
				if err == nil {
					global.LbLogger.Info(fmt.Sprintf("Successfully connected to nsqlookupd."))
					c.isStarted = true
					break
				}
				global.LbLogger.Info(fmt.Sprintf("Failed to connect to nsqlookupd: %v. Retrying in 5 seconds...", err))
				time.Sleep(5 * time.Second) // 连接失败后等待 5 秒再重试
			}
		})
	}()
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
	cfg.LookupdPollInterval = 10 * time.Second
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
	c.NSQConsumer.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		if len(m.Body) == 0 {
			// Returning nil will automatically send a FIN command to NSQ to mark the message as processed.
			// In this case, a message with an empty body is simply ignored/discarded.
			return nil
		}
		c.Logger.Info(string(m.Body))
		return nil
	},
	))

	return nil
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
	return c, nil
}
