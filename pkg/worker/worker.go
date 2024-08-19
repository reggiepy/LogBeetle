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
func NewWorker(options ...Option) (*Worker, error) {
	w := &Worker{}
	for _, opt := range options {
		if err := opt.apply(w); err != nil {
			return nil, err
		}
	}
	return w, nil
}


// Run 启动工作器，执行具体的工作
func (w *Worker) Run() {
	if w.Wg != nil {
		w.Wg.Add(1)
	}

	go func() {
		defer func() {
			if w.Wg != nil {
				w.Wg.Done()
			}
		}()
		fmt.Printf("【%s】Working Start...\n", w.Name)
		if w.Start != nil {
			w.Start()
		}

		<-w.Ctx.Done()

		fmt.Printf("【%s】Worker stopped\n", w.Name)
		if w.Stop != nil {
			w.Stop()
		}
	}()
}
