package nsq_consumer

import (
	"fmt"
	"path"
	"sync"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/reggiepy/LogBeetle/global"
	"github.com/reggiepy/LogBeetle/pkg/consumer"
	"github.com/reggiepy/goutils/logutil/zapLogger"
	"go.uber.org/zap"
)

// NSQLogConsumer 表示一个 NSQ 日志消费实例
type NSQLogConsumer struct {
	name       string         // 消费者名称
	logger     *zap.Logger    // 系统级 logger
	isStarted  bool           // 是否已连接成功
	startOnce  sync.Once      // 确保连接只启动一次
	stopOnce   sync.Once      // 确保停止逻辑只执行一次
	nsqTopic   string         // NSQ topic
	nsqAddress string         // NSQD 地址

	nsqAuthSecret        string         // NSQ 鉴权密钥（可选）
	nsqConsumer          *nsq.Consumer  // NSQ 消费者对象
	handlerLogger        *zap.Logger    // 用于写日志到文件的 logger
	handlerFileName      string         // 写入日志的文件名
	handlerLoggerCleanup func()         // 日志 logger 清理函数
}

// GetName 返回消费者名称
func (c *NSQLogConsumer) GetName() string {
	return c.name
}

// GetType 返回消费者类型
func (c *NSQLogConsumer) GetType() string {
	return consumer.ChannelConsumer.String()
}

// Stop 停止 NSQ 消费者和日志记录器
func (c *NSQLogConsumer) Stop() error {
	c.stopOnce.Do(func() {
		if c.nsqConsumer != nil {
			c.nsqConsumer.Stop()
			c.logger.Info(fmt.Sprintf("Stopped NSQ consumer [%s]", c.name))
		}
		if c.handlerLoggerCleanup != nil {
			c.handlerLoggerCleanup()
		}
	})
	return nil
}

// Start 启动消费者，连接 NSQD（失败则自动重试）
func (c *NSQLogConsumer) Start() error {
	go c.startOnce.Do(func() {
		for {
			if err := c.nsqConsumer.ConnectToNSQD(c.nsqAddress); err == nil {
				c.logger.Info(fmt.Sprintf("Connected to NSQ [%s]", c.nsqAddress))
				c.isStarted = true
				break
			} else {
				c.logger.Warn(fmt.Sprintf("Failed to connect to NSQ [%s]: %v, retrying...", c.nsqAddress, err))
				time.Sleep(5 * time.Second)
			}
		}
	})
	return nil
}

// initializeNSQConsumer 初始化 NSQ 消费者对象及处理函数
func (c *NSQLogConsumer) initializeNSQConsumer() error {
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = 10 * time.Second

	// 设置 NSQ 鉴权
	if c.nsqAuthSecret != "" {
		if err := cfg.Set("auth_secret", c.nsqAuthSecret); err != nil {
			return fmt.Errorf("set auth_secret failed: %w", err)
		}
	}

	// 创建 NSQ 消费者
	nsqConsumer, err := nsq.NewConsumer(c.nsqTopic, "consumer", cfg)
	if err != nil {
		return fmt.Errorf("create NSQ consumer failed: %w", err)
	}

	// 设置 NSQ 消费日志记录器
	nsqConsumer.SetLogger(&zapNSQLogger{logger: c.logger.Named("nsq")}, nsq.LogLevelInfo)

	// 设置消息处理函数
	nsqConsumer.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
		if len(m.Body) == 0 {
			return nil // 忽略空消息
		}
		c.handlerLogger.Info(string(m.Body))
		return nil
	}))

	c.nsqConsumer = nsqConsumer
	return nil
}

// initializeLogger 初始化日志记录器并绑定清理函数
func (c *NSQLogConsumer) initializeLogger() error {
	logFile := path.Join(global.LbConfig.ConsumerConfig.LogPath, c.handlerFileName)

	logCfg := zapLogger.NewLoggerConfig(
		zapLogger.WithFile(logFile),
		zapLogger.WithLogFormat("logfmt"),
	)

	logger, cleanup := zapLogger.NewLogger(logCfg)
	c.handlerLogger = logger
	c.handlerLoggerCleanup = cleanup
	return nil
}

// NewNSQLogConsumer 构造并初始化一个新的 NSQ 日志消费者实例
func NewNSQLogConsumer(name, nsqAddress string, opts ...Options) (*NSQLogConsumer, error) {
	c := &NSQLogConsumer{
		name:       name,
		nsqAddress: nsqAddress,
	}

	// 应用可选配置项
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	// 默认值补全
	if c.nsqTopic == "" {
		c.nsqTopic = name
	}
	if c.handlerFileName == "" {
		c.handlerFileName = name + ".log"
	}
	if c.logger == nil {
		c.logger = global.LbLogger.Named("consumer")
	}

	// 初始化日志记录器
	if err := c.initializeLogger(); err != nil {
		return nil, err
	}

	// 初始化 NSQ 消费者
	if err := c.initializeNSQConsumer(); err != nil {
		return nil, err
	}

	return c, nil
}
