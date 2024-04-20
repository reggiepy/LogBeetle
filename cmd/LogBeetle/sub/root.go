package sub

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/reggiepy/LogBeetle/pkg/config"
	"github.com/reggiepy/LogBeetle/pkg/consumer/logconsumer"
	"github.com/reggiepy/LogBeetle/pkg/consumer/nsqconsumer"
	"github.com/reggiepy/LogBeetle/pkg/convert"
	"github.com/reggiepy/LogBeetle/pkg/logger"
	"github.com/reggiepy/LogBeetle/pkg/producer/nsqproducer"
	"github.com/reggiepy/LogBeetle/pkg/worker"
	"github.com/reggiepy/LogBeetle/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	showVersion bool
)

func init() {
	// 使用Viper加载环境变量
	viper.SetEnvPrefix("lb") // 设置环境变量前缀
	viper.SetDefault("config", "./log-beetle.yaml")
	viper.AutomaticEnv() // 自动加载环境变量

	// 初始化配置
	cobra.OnInitialize(initConfig)

	// 设置全局标志
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "show version information")

	// 添加命令行参数
	rootCmd.Flags().StringP("config", "c", "", "config file")
	rootCmd.Flags().String("log-file", "", "file to log")
	rootCmd.Flags().String("consumer-log-path", "", "Path to consumer log")
	rootCmd.Flags().String("nsq-address", "", "Address of NSQ")
	rootCmd.Flags().String("nsq-auth-secret", "", "NSQ auth secret")

	// 绑定命令行参数到Viper
	_ = viper.BindPFlag("log_file", rootCmd.Flags().Lookup("log-file"))
	_ = viper.BindPFlag("nsq_address", rootCmd.Flags().Lookup("nsq-address"))
	_ = viper.BindPFlag("nsq_auth_secret", rootCmd.Flags().Lookup("nsq-auth-secret"))
	_ = viper.BindPFlag("consumer_log_path", rootCmd.Flags().Lookup("consumer-log-path"))
	_ = viper.BindPFlag("config", rootCmd.Flags().Lookup("config"))
}

var rootCmd = cobra.Command{
	Use:     "LogBeetle",
	Short:   "LogBeetle help",
	Long:    `LogBeetle help`,
	Version: "",
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = config.ShowConfig("simple")
		loggerConfig, err := convert.ConfigToLoggerConfig(config.Instance)
		err = logger.InitLogger(*loggerConfig)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		serverStart()
	},
}

func initConfig() {
	// 从Viper获取配置文件路径
	configFile := viper.GetString("config")
	// 初始化配置
	config.Init(configFile)

	// 从Viper获取日志文件路径
	logPath := viper.GetString("log_file")
	// 如果日志文件路径不为空，则更新配置中的日志文件路径
	if logPath != "" {
		config.Instance.LogConfig.LogFile = logPath
	}

	// 从Viper获取消费者日志路径
	consumerLogPath := viper.GetString("consumer_log_path")
	// 如果消费者日志路径不为空，则更新配置中的消费者日志路径
	if consumerLogPath != "" {
		config.Instance.ConsumerConfig.LogPath = consumerLogPath
	}

	// 从Viper获取NSQ地址
	nsqAddress := viper.GetString("nsq_address")
	// 如果NSQ地址不为空，则更新配置中的NSQ地址
	if nsqAddress != "" {
		config.Instance.NSQConfig.NSQDAddress = nsqAddress
	}

	// 从Viper获取NSQ认证密钥
	nsqAuthSecret := viper.GetString("nsq_auth_secret")
	// 如果NSQ认证密钥不为空，则更新配置中的NSQ认证密钥
	if nsqAuthSecret != "" {
		config.Instance.NSQConfig.AuthSecret = nsqAuthSecret
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func serverStart() {
	// 创建一个上下文，以便能够在主程序退出时取消所有 Goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 使用 WaitGroup 等待所有 Goroutine 结束
	var wg sync.WaitGroup

	// 设置路由
	router := web.SetupRouter()

	// 获取配置
	nsqConfig := config.Instance.NSQConfig
	consumerConfig := config.Instance.ConsumerConfig
	if len(consumerConfig.NSQConsumers) == 0 {
		_ = fmt.Errorf("consumer config is empty")
	}

	// 定义工作线程
	workers := []*worker.Worker{
		// Web 服务器
		worker.NewWorker(
			worker.WithName("webserver"),
			worker.WithCtx(ctx),
			worker.WithWg(&wg),
			worker.WithStop(func() {
			}),
			worker.WithStart(func() {
				go func() {
					err := router.Run(":1233")
					if err != nil {
						fmt.Printf("Error: %v\n", err)
					}
				}()
			}),
		),
		// NSQ 生产者和消费者
		worker.NewWorker(
			worker.WithName("nsqProducer"),
			worker.WithCtx(ctx),
			worker.WithWg(&wg),
			worker.WithStop(func() {
				nsqproducer.StopProducer()
				logconsumer.StopConsumers()
			}),
			worker.WithStart(func() {
				// 初始化 NSQ 生产者
				nsqproducer.InitProducer(nsqproducer.ProducerConfig{
					Address:    nsqConfig.NSQDAddress,
					AuthSecret: nsqConfig.AuthSecret,
				})

				nc := nsqconsumer.NewNsqConsumer(
					"test",
					nsqConfig.NSQDAddress,
					nsqconsumer.WithAuthSecret(nsqConfig.AuthSecret),
				)
				c := logconsumer.NewLogConsumer(
					"test",
					"test.log",
					nc,
				)
				// 添加日志消费者
				logconsumer.AddConsumer(c)

				// 添加其他消费者
				for _, consumerConfig := range consumerConfig.NSQConsumers {
					nc := nsqconsumer.NewNsqConsumer(
						consumerConfig.Topic,
						nsqConfig.NSQDAddress,
						nsqconsumer.WithAuthSecret(nsqConfig.AuthSecret),
					)
					c := logconsumer.NewLogConsumer(
						consumerConfig.Name,
						consumerConfig.FileName,
						nc,
					)
					// 添加日志消费者
					logconsumer.AddConsumer(c)
				}
			}),
		),
	}

	// 启动工作线程
	for _, w := range workers {
		wg.Add(1)
		w.Run()
	}

	// 捕获信号，以优雅地退出程序
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// 取消所有 Goroutine
	cancel()
	wg.Wait()

	fmt.Println("Main program stopped")
}
