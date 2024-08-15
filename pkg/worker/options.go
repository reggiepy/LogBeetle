package worker

import (
	"context"
	"sync"
)

// Option 定义了一个接口，表示一个配置选项
type Option interface {
	apply(*Worker) error
}

// OptionFunc 是一种为函数类型定义的新类型，它的参数需要传递指针
type OptionFunc func(w *Worker) error

// apply 方法将 OptionFunc 类型转换为 Option 接口类型
func (f OptionFunc) apply(w *Worker) error {
	return f(w)
}

// WithName 设置工作器的名称选项
func WithName(name string) Option {
	return OptionFunc(func(w *Worker) error {
		w.Name = name
		return nil
	})
}

// WithStop 设置工作器的停止函数选项
func WithStop(stop func()) Option {
	return OptionFunc(func(w *Worker) error {
		w.Stop = stop
		return nil
	})
}

// WithStart 设置工作器的启动函数选项
func WithStart(start func()) Option {
	return OptionFunc(func(w *Worker) error {
		w.Start = start
		return nil
	})
}

// WithCtx 设置工作器的上下文选项
func WithCtx(ctx context.Context) Option {
	return OptionFunc(func(w *Worker) error {
		w.Ctx = ctx
		return nil
	})
}

// WithWg 设置工作器的等待组选项
func WithWg(wg *sync.WaitGroup) Option {
	return OptionFunc(func(w *Worker) error {
		w.Wg = wg
		return nil
	})
}
