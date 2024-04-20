package logconsumer

import (
	"fmt"
	"path"

	"github.com/natefinch/lumberjack"
	"github.com/reggiepy/LogBeetle/pkg/config"
	"github.com/reggiepy/LogBeetle/pkg/consumer/nsqconsumer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConsumer struct {
	Name        string
	LogFileName string
	LogFile     *lumberjack.Logger
	// Logger    zerolog.Logger
	Logger2  *zap.Logger
	Consumer *nsqconsumer.NsqConsumer
}

func (c *LogConsumer) Close() {
	var err error
	if c.LogFile != nil {
		fmt.Printf("关闭 【%s】 Consumer LogFile\n", c.Name)
		err = c.Logger2.Sync()
		if err != nil {
			fmt.Printf("Sync 【%s】  Logger2 失败: %v", c.Name, err)
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

type Options func(logConsumer *LogConsumer) error

// WithLogFile 设置日志文件
func WithLogFile(logFile *lumberjack.Logger) Options {
	return func(logConsumer *LogConsumer) error {
		logConsumer.LogFile = logFile
		return nil
	}
}

func NewLogConsumer(name string, logFileName string, consumer *nsqconsumer.NsqConsumer, opts ...Options) *LogConsumer {
	c := &LogConsumer{
		Name:        name,
		LogFileName: logFileName,
		Consumer:    consumer,
	}
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			panic((err).(any))
		}
	}

	// 创建 lumberjack.Logger 实例用于日志切割
	consumerConfig := config.Instance.ConsumerConfig
	filePath := path.Join(consumerConfig.LogPath, c.LogFileName)
	c.LogFile = &lumberjack.Logger{
		Filename:   filePath, // 日志文件名
		MaxSize:    1,        // 日志文件大小限制，单位为 MB
		MaxBackups: 5,        // 最大保留的旧日志文件数量
		MaxAge:     30,       // 旧日志文件保留天数
		Compress:   false,    // 是否压缩旧日志文件
	}
	//// 创建 ConsoleWriter 以确保在 Windows 平台下的正确编码
	//writer := zerolog.ConsoleWriter{
	//	Out:        logFile,
	//	TimeFormat: time.RFC3339,
	//	NoColor:    true,
	//}
	//logger := log.Output(writer)

	// 修改时间编码器，在日志文件中使用大写字母记录日志级别
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 创建Zap核心
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(c.LogFile),
		zapcore.InfoLevel,
	)

	// 创建Logger
	c.Logger2 = zap.New(core)
	c.Consumer.AddHandler(&nsqconsumer.MessageHandler{
		Handler: func(message []byte) error {
			c.Logger2.Info(string(message))
			return nil
		},
	})
	err := c.Consumer.Connect()
	if err != nil {
		panic((err).(any))
	}
	return c
}
