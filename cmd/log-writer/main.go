package main

import (
	"context"
	"fmt"
	"github.com/reggiepy/LogBeetle/nsqworker"
	"github.com/reggiepy/LogBeetle/web"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Worker struct {
	Name string `json:"name"`
	Stop func()
	Run  func()
	ctx  context.Context
	wg   *sync.WaitGroup
}

func NewWorker() *Worker {
	return &Worker{}
}

func (w *Worker) SetName(name string) {
	w.Name = name
}

func (w *Worker) SetRun(run func()) {
	w.Run = run
}

func (w *Worker) SetStop(stop func()) {
	w.Stop = stop
}

func (w *Worker) Start() {
	go func() {
		defer w.wg.Done()
		for {
			select {
			case <-w.ctx.Done():
				fmt.Printf("【%s】Worker stopped\n", w.Name)
				w.Stop()
				return
			default:
				// Do some work here
			}
		}
	}()
	go func() {
		fmt.Printf("【%s】Working Start...\n", w.Name)
		w.Run()
	}()
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
	authSecret := "%n&yFA2JD85z^g"
	address := "192.168.1.110:4150"
	workers := []*Worker{
		&Worker{
			ctx:  ctx,
			wg:   &wg,
			Name: "webserver",
			Stop: func() {

			},
			Run: func() {
				go func() {
					router := web.SetupRouter()
					err := router.Run(":1233")
					if err != nil {
						fmt.Printf("Error: %v\n", err)
					}
				}()
			},
		},
		&Worker{
			ctx:  ctx,
			wg:   &wg,
			Name: "nsqProducer",
			Stop: func() {
				nsqworker.StopConsumers()
				nsqworker.StopProducer()
			},
			Run: func() {
				nsqworker.InitProducer(nsqworker.ProducerConfig{
					Address:    address,
					AuthSecret: authSecret,
				})
				nsqworker.AddConsumer(nsqworker.NewConsumer(nsqworker.ConsumerConfig{
					Address:    address,
					AuthSecret: authSecret,
					Topic:      "test",
					Channel:    "test_channel",
					Handler: &nsqworker.MessageHandler{Handler: func(message []byte) error {
						fmt.Println(string(message))
						return nil
					}},
				}))
			},
		},
	}
	for _, worker := range workers {
		wg.Add(1)
		worker.Start()
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
