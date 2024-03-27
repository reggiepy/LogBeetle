package sub

import (
	"context"
	"fmt"
	"github.com/reggiepy/LogBeetle/pkg/config"
	"github.com/reggiepy/LogBeetle/pkg/consumer/logconsumer"
	"github.com/reggiepy/LogBeetle/pkg/consumer/nsqconsumer"
	"github.com/reggiepy/LogBeetle/pkg/logger"
	"github.com/reggiepy/LogBeetle/pkg/producer/nsqproducer"
	"github.com/reggiepy/LogBeetle/pkg/worker"
	"github.com/reggiepy/LogBeetle/web"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	//配置文件
	configFile  string
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
	Run: func(cmd *cobra.Command, args []string) {
		_ = config.ShowConfig("simple")
		err := logger.InitLogger(config.Instance)
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

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func serverStart() {
	// 创建一个上下文，以便能够在主程序退出时取消所有 Goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 使用 WaitGroup 等待所有 Goroutine 结束
	var wg sync.WaitGroup
	router := web.SetupRouter()

	nsqConfig := config.Instance.NSQConfig
	consumerConfig := config.Instance.ConsumerConfig
	workers := []*worker.Worker{
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
		worker.NewWorker(
			worker.WithName("nsqProducer"),
			worker.WithCtx(ctx),
			worker.WithWg(&wg),
			worker.WithStop(func() {
				nsqproducer.StopProducer()
				logconsumer.StopConsumers()
			}),
			worker.WithStart(func() {
				nsqproducer.InitProducer(nsqproducer.ProducerConfig{
					Address:    nsqConfig.NSQDAddress,
					AuthSecret: nsqConfig.AuthSecret,
				})
				logconsumer.AddConsumer(
					logconsumer.NewLogConsumer(
						"test",
						nsqconsumer.NsqConsumerConfig{
							Address:    nsqConfig.NSQDAddress,
							AuthSecret: nsqConfig.AuthSecret,
							Topic:      "test",
							Channel:    "test_channel",
						}, "test.log"),
				)
				for _, consumerConfig := range consumerConfig.Consumers {
					fmt.Println(consumerConfig)
					logconsumer.AddConsumer(
						logconsumer.NewLogConsumer(
							consumerConfig.Name,
							nsqconsumer.NsqConsumerConfig{
								Address:    nsqConfig.NSQDAddress,
								AuthSecret: nsqConfig.AuthSecret,
								Topic:      consumerConfig.Topic,
								Channel:    consumerConfig.Channel,
							}, consumerConfig.FileName),
					)
				}
			}),
		),
	}
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
