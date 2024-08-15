package worker

import (
	"context"
	"fmt"
	"sync"
)

// Worker 包含了工作器的配置选项
type Worker struct {
	Name  string          `json:"name"` // 使用一致的命名风格
	Start func()          // 启动函数
	Stop  func()          // 停止函数
	Ctx   context.Context // 上下文
	Wg    *sync.WaitGroup // 等待组
}

// NewWorker 创建一个新的 Worker 实例，接受一个或多个选项
func NewWorker(options ...Option) *Worker {
	var (
		err error
		w   = &Worker{}
	)
	for _, opt := range options {
		err = opt.apply(w)
		if err != nil {
			panic(err.(any))
		}
	}
	return w
}

// Run 启动工作器，执行具体的工作
func (w *Worker) Run() {
	go func() {
		defer w.Wg.Done()
		for {
			select {
			case <-w.Ctx.Done():
				fmt.Printf("【%s】Worker stopped\n", w.Name)
				w.Stop()
				return
				// default:
				//	// 在这里执行具体的工作
			}
		}
	}()
	go func() {
		fmt.Printf("【%s】Working Start...\n", w.Name)
		w.Start()
	}()
}
