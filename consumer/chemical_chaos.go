package consumer

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/nsqio/go-nsq"
	"github.com/reggiepy/LogBeetle/nsqworker"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"time"
)

type LogConsumer struct {
	Name      string
	NsqConfig nsqworker.ConsumerConfig
	LogFile   *lumberjack.Logger
	Writer    io.Writer
	Logger    zerolog.Logger
	Consumer  *nsq.Consumer
}

func NewLogConsumer(name string, nsqConfig nsqworker.ConsumerConfig, fileName string) *LogConsumer {
	// 创建 lumberjack.Logger 实例用于日志切割
	logFile := &lumberjack.Logger{
		Filename:   fileName, // 日志文件名
		MaxSize:    10,       // 日志文件大小限制，单位为 MB
		MaxBackups: 5,        // 最大保留的旧日志文件数量
		MaxAge:     30,       // 旧日志文件保留天数
		Compress:   true,     // 是否压缩旧日志文件
	}
	// 创建 ConsoleWriter 以确保在 Windows 平台下的正确编码
	writer := zerolog.ConsoleWriter{
		Out:        logFile,
		TimeFormat: time.RFC3339,
		NoColor:    true,
	}
	logger := log.Output(writer)

	c := &LogConsumer{
		Name:    name,
		LogFile: logFile,
		Writer:  &writer,
		Logger:  logger,
	}
	nsqConfig.Handler = &nsqworker.MessageHandler{
		Handler: c.Handler,
	}
	consumer := nsqworker.NewConsumer(nsqConfig)
	c.Consumer = consumer
	return c
}

func (c *LogConsumer) Handler(message []byte) error {
	c.Logger.Info().Msg(string(message))
	return nil
}
func (c *LogConsumer) Close() {
	if c.LogFile != nil {
		fmt.Printf("关闭 【%s】 Consumer LogFile\n", c.Name)
		_ = c.LogFile.Close()
	}
	if c.Consumer != nil {
		fmt.Printf("关闭 【%s】 Consumer Consumer\n", c.Name)
		c.Consumer.Stop()
	}
}
