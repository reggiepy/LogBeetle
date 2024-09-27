package signailUtils

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	shutdownHooks      []func()         // 退出时需要执行的钩子函数列表
	exitMessageHandler func(msg string) // 退出消息处理函数
	mu                 sync.Mutex       // 确保并发情况下对 shutdownHooks 的操作安全
	once               sync.Once        // 确保退出函数只执行一次
	wg                 sync.WaitGroup   // 确保所有协程都在程序退出前完成
)

// SetExitMessageHandler - 设置自定义的退出消息处理函数
func SetExitMessageHandler(handler func(msg string)) {
	exitMessageHandler = handler // 将自定义处理函数赋值
}

// GetExitMessageHandler - 获取退出消息处理函数，默认为标准输出
func GetExitMessageHandler() func(msg string) {
	if exitMessageHandler != nil {
		return exitMessageHandler // 返回自定义的处理函数
	}
	// 默认处理函数，打印消息到控制台
	return func(msg string) {
		fmt.Println(msg)
	}
}

// OnExit - 注册退出处理函数
func OnExit(fnExit func()) {
	mu.Lock()                                     // 获取锁以确保安全
	defer mu.Unlock()                             // 解锁
	shutdownHooks = append(shutdownHooks, fnExit) // 将退出处理函数添加到列表
	wg.Add(1)                                     // 为每个退出钩子增加一个等待组
}

// ExecuteShutdownHooks - 执行所有注册的退出钩子函数，只执行一次
func ExecuteShutdownHooks() {
	GetExitMessageHandler()("准备退出,执行清理操作")
	once.Do(func() { // 确保只执行一次
		mu.Lock()         // 获取锁以确保安全
		defer mu.Unlock() // 解锁
		for _, fnExit := range shutdownHooks {
			go func(fn func()) {
				defer wg.Done() // 每个钩子执行完毕后调用 Done
				fn()            // 执行退出处理函数
			}(fnExit)
		}
	})
	GetExitMessageHandler()("程序清理完毕。")
}

// WaitExit - 同步等待退出信号后退出，可指定自定义退出处理函数
func WaitExit(timeout time.Duration) {
	osc := make(chan os.Signal, 1)                      // 创建信号通道
	signal.Notify(osc, syscall.SIGTERM, syscall.SIGINT) // 监听 SIGTERM 和 SIGINT 信号

	<-osc // 阻塞等待信号
	ExecuteShutdownHooks() // 执行已注册的退出钩子函数

	// 等待所有退出钩子完成或超时
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done) // 退出钩子执行完，关闭通道
	}()

	select {
	case <-done: // 等待完成
		GetExitMessageHandler()("所有清理已完成，程序正常退出。")
	case <-time.After(timeout): // 设置超时时间
		GetExitMessageHandler()("退出处理因超时而中止。")
	}
}
