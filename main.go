package main

import (
	"context"
	"fmt"
	"github.com/reggiepy/LogBeetle/consumer"
	"github.com/reggiepy/LogBeetle/nsqworker"
	"github.com/reggiepy/LogBeetle/web"
	"github.com/reggiepy/LogBeetle/worker"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

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

	authSecret := "%n&yFA2JD85z^g"
	address := "192.168.1.110:4150"
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
				nsqworker.StopProducer()
				consumer.StopConsumers()
			}),
			worker.WithStart(func() {
				nsqworker.InitProducer(nsqworker.ProducerConfig{
					Address:    address,
					AuthSecret: authSecret,
				})
				consumer.AddConsumer(
					consumer.NewLogConsumer(
						"chemical_chaos",
						nsqworker.ConsumerConfig{
							Address:    address,
							AuthSecret: authSecret,
							Topic:      "test",
							Channel:    "test_channel",
						}, "chemical_chaos.log"),
				)
				consumer.AddConsumer(
					consumer.NewLogConsumer(
						"chemical_db",
						nsqworker.ConsumerConfig{
							Address:    address,
							AuthSecret: authSecret,
							Topic:      "chemical-db",
							Channel:    "test_channel",
						}, "chemical_db.log"),
				)
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