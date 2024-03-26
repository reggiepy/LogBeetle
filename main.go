package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/reggiepy/LogBeetle/pkg/config"
	"github.com/reggiepy/LogBeetle/pkg/consumer/logconsumer"
	"github.com/reggiepy/LogBeetle/pkg/consumer/nsqconsumer"
	"github.com/reggiepy/LogBeetle/pkg/logger"
	"github.com/reggiepy/LogBeetle/pkg/producer/nsqproducer"
	"github.com/reggiepy/LogBeetle/pkg/worker"
	"github.com/reggiepy/LogBeetle/web"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var configFile = flag.String("config", "./log-beetle.yaml", "配置文件路径")

func init() {
	flag.Parse()

	// 初始化配置
	conf := config.Init(*configFile)

	err := logger.InitLogger(conf)
	if err != nil {
		fmt.Println(err.Error())
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
func main() {
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
