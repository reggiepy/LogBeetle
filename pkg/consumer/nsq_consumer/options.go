package nsq_consumer

import (
	"go.uber.org/zap"
)

type Options func(n *NSQLogConsumer) error

// WithAuthSecret 设置NSQ消费者的认证秘钥
func WithNSQAuthSecret(authSecret string) Options {
	return func(n *NSQLogConsumer) error {
		n.nsqAuthSecret = authSecret
		return nil
	}
}

// WithLogFileName 设置NSQ消费者的日志文件名
func WithLogFileName(logFileName string) Options {
	return func(n *NSQLogConsumer) error {
		n.handlerFileName = logFileName
		return nil
	}
}

// WithLogger 设置NSQ消费者的日志记录器
func WithLogger(logger *zap.Logger) Options {
	return func(n *NSQLogConsumer) error {
		n.logger = logger
		return nil
	}
}

// WithNSQTopic 设置NSQ消费者的主题
func WithNSQTopic(topic string) Options {
	return func(n *NSQLogConsumer) error {
		n.nsqTopic = topic
		return nil
	}
}
