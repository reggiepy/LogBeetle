package consumer

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/nsqio/go-nsq"
	"github.com/reggiepy/LogBeetle/pkg/config"
	"github.com/reggiepy/LogBeetle/pkg/nsqworker"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path"
)

type LogConsumer struct {
	Name      string
	NsqConfig nsqworker.ConsumerConfig
	LogFile   *lumberjack.Logger
	//Logger    zerolog.Logger
	Logger2  *zap.Logger
	Consumer *nsq.Consumer
}

func NewLogConsumer(name string, nsqConfig nsqworker.ConsumerConfig, fileName string) *LogConsumer {
	// 创建 lumberjack.Logger 实例用于日志切割
	consumerConfig := config.Instance.ConsumerConfig
	filePath := path.Join(consumerConfig.LogPath, fileName)
	logFile := &lumberjack.Logger{
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

	//修改时间编码器，在日志文件中使用大写字母记录日志级别
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 创建Zap核心
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(logFile),
		zapcore.InfoLevel,
	)

	// 创建Logger
	logger2 := zap.New(core)

	c := &LogConsumer{
		Name:    name,
		LogFile: logFile,
		//Logger:  logger,
		Logger2: logger2,
	}
	nsqConfig.Handler = &nsqworker.MessageHandler{
		Handler: c.Handler,
	}
	consumer := nsqworker.NewConsumer(nsqConfig)
	c.Consumer = consumer
	return c
}

func (c *LogConsumer) Handler(message []byte) error {
	c.Logger2.Info(string(message))
	return nil
}
func (c *LogConsumer) Close() {
	var err error
	if c.LogFile != nil {
		fmt.Printf("关闭 【%s】 Consumer LogFile\n", c.Name)
		err = c.Logger2.Sync()
		if err != nil {
			fmt.Printf("Sync Logger2 失败: %v", err)
		}
		err = c.LogFile.Close()
		if err != nil {
			fmt.Printf("关闭日志文件失败: %v", err)
		}
	}
	if c.Consumer != nil {
		fmt.Printf("关闭 【%s】 Consumer Consumer\n", c.Name)
		c.Consumer.Stop()
	}
}
