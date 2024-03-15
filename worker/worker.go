package worker

import (
	"context"
	"fmt"
	"sync"
)

// Worker 包含了工作器的配置选项
type Worker struct {
	Options *Options `json:"options"` // 使用更加一致的名称，表示它可以包含多个选项
}

// NewWorker 创建一个新的 Worker 实例，接受一个或多个选项
func NewWorker(options ...Option) *Worker {
	var err error
	Options := new(Options) // 使用一致的命名风格
	for _, opt := range options {
		err = opt.apply(Options)
		if err != nil {
			panic(err)
		}
	}
	return &Worker{
		Options: Options,
	}
}

// Run 启动工作器，执行具体的工作
func (w *Worker) Run() {
	go func() {
		defer w.Options.Wg.Done()
		for {
			select {
			case <-w.Options.Ctx.Done():
				fmt.Printf("【%s】Worker stopped\n", w.Options.Name)
				w.Options.Stop()
				return
			default:
				// 在这里执行具体的工作
			}
		}
	}()
	go func() {
		fmt.Printf("【%s】Working Start...\n", w.Options.Name)
		w.Options.Start()
	}()
}

// Options 包含了工作器的所有配置选项
type Options struct {
	Name  string          `json:"name"` // 使用一致的命名风格
	Start func()          // 启动函数
	Stop  func()          // 停止函数
	Ctx   context.Context // 上下文
	Wg    *sync.WaitGroup // 等待组
}

// Option 定义了一个接口，表示一个配置选项
type Option interface {
	apply(option *Options) error
}

// OptionFunc 是一种为函数类型定义的新类型，它的参数需要传递指针
type OptionFunc func(option *Options) error

// apply 方法将 OptionFunc 类型转换为 Option 接口类型
func (f OptionFunc) apply(opt *Options) error {
	return f(opt)
}

// WithName 设置工作器的名称选项
func WithName(name string) Option {
	return OptionFunc(func(option *Options) error {
		option.Name = name
		return nil
	})
}

// WithStop 设置工作器的停止函数选项
func WithStop(stop func()) Option {
	return OptionFunc(func(option *Options) error {
		option.Stop = stop
		return nil
	})
}

// WithStart 设置工作器的启动函数选项
func WithStart(start func()) Option {
	return OptionFunc(func(option *Options) error {
		option.Start = start
		return nil
	})
}

// WithCtx 设置工作器的上下文选项
func WithCtx(ctx context.Context) Option {
	return OptionFunc(func(option *Options) error {
		option.Ctx = ctx
		return nil
	})
}

// WithWg 设置工作器的等待组选项
func WithWg(wg *sync.WaitGroup) Option {
	return OptionFunc(func(option *Options) error {
		option.Wg = wg
		return nil
	})
}
