package consumer

import (
	"go.uber.org/zap"
	"io"
)

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
