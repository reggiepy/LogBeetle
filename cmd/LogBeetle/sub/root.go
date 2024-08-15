package sub

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/reggiepy/LogBeetle/boot"
	"github.com/reggiepy/LogBeetle/global"
	"github.com/reggiepy/LogBeetle/goutils/enumUtils"
	"github.com/reggiepy/LogBeetle/version"

	"github.com/reggiepy/LogBeetle/pkg/consumer"
	"github.com/reggiepy/LogBeetle/pkg/producer"

	"github.com/reggiepy/LogBeetle/pkg/worker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	showVersion  bool
	configFormat = enumUtils.NewEnum([]string{"humanReadable", "simple"}, "humanReadable")
)

func init() {
	viper.SetDefault("config", "./log-beetle.yaml")
	// 设置全局标志
	rootCmd.PersistentFlags().BoolVarP(&showVersion, "version", "v", false, "show version information")
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file")

	// 添加命令行参数
	rootCmd.Flags().String("log-file", "", "file to log")
	rootCmd.Flags().String("consumer-log-path", "", "Path to consumer log")
	rootCmd.Flags().String("nsq-address", "", "Address of NSQ")
	rootCmd.Flags().String("nsq-auth-secret", "", "NSQ auth secret")

	// 绑定命令行参数到Viper
	_ = viper.BindPFlag("log_file", rootCmd.Flags().Lookup("log-file"))
	_ = viper.BindPFlag("nsq_address", rootCmd.Flags().Lookup("nsq-address"))
	_ = viper.BindPFlag("nsq_auth_secret", rootCmd.Flags().Lookup("nsq-auth-secret"))
	_ = viper.BindPFlag("consumer_log_path", rootCmd.Flags().Lookup("consumer-log-path"))
	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
}

var rootCmd = cobra.Command{
	Use:   "LogBeetle",
	Short: "LogBeetle help",
	Long:  `LogBeetle help`,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if showVersion {
			fmt.Println(version.Full())
			return nil
		}
		global.LbViper = boot.Viper()
		global.LbLogger = boot.Log()
		boot.Boot()
		StartServer()
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func StartServer() {
	// 创建一个上下文，以便能够在主程序退出时取消所有 Goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 使用 WaitGroup 等待所有 Goroutine 结束
	var wg sync.WaitGroup
	// 获取配置
	nsqConfig := global.LbConfig.NSQConfig
	consumerConfig := global.LbConfig.ConsumerConfig
	if len(consumerConfig.NSQConsumers) == 0 {
		_ = fmt.Errorf("consumer config is empty")
	}

	consumerManager := consumer.NewLogBeetleConsumerManager()
	producerManager := producer.NewLogBeetleProducerManager()

	// 定义工作线程
	workers := []*worker.Worker{
		// NSQ 生产者和消费者
		worker.NewWorker(
			worker.WithName("nsq worker"),
			worker.WithCtx(ctx),
			worker.WithWg(&wg),
			worker.WithStop(func() {
				producerManager.Stop()
				consumerManager.Stop()
			}),
			worker.WithStart(func() {
				// 初始化 NSQ 生产者
				p, err := producer.NewNSQProducer(
					producer.WithNSQAddress(nsqConfig.NSQDAddress),
					producer.WithNSQAuthSecret(nsqConfig.AuthSecret),
				)
				if err != nil {
					fmt.Printf("error creating producer: %v\n", err)
				} else {
					producerManager.Add(p)
					producer.InitInstance(p)
				}

				c, err := consumer.NewNSQLogConsumer(
					consumer.WithName("test"),
					consumer.WithLogFileName("test.log"),
					consumer.WithNSQTopic("test"),
					consumer.WithNSQAddress(nsqConfig.NSQDAddress),
					consumer.WithNSQAuthSecret(nsqConfig.AuthSecret),
				)
				if err != nil {
					fmt.Printf("error creating consumer %s: %v\n", "test", err)
				} else {
					consumerManager.Add(c)
					fmt.Printf("consumer %s added\n", "test")
					global.LbLogger.Info(fmt.Sprintf("consumer %s added\n", "test"))
				}

				// 添加其他消费者
				for _, consumerConfig := range consumerConfig.NSQConsumers {
					c, err := consumer.NewNSQLogConsumer(
						consumer.WithName(consumerConfig.Name),
						consumer.WithLogFileName(consumerConfig.FileName),
						consumer.WithNSQTopic(consumerConfig.Topic),
						consumer.WithNSQAddress(nsqConfig.NSQDAddress),
						consumer.WithNSQAuthSecret(nsqConfig.AuthSecret),
					)
					if err != nil {
						fmt.Printf("error creating consumer %s: %v\n", consumerConfig.Topic, err)
					} else {
						consumerManager.Add(c)
						fmt.Printf("consumer %s added\n", consumerConfig.Name)
					}
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
