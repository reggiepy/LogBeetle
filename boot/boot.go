package boot

import (
	"context"
	"errors"
	"fmt"
	"github.com/reggiepy/LogBeetle/global"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Boot() {
	addr := fmt.Sprintf("%s:%d", "127.0.0.1", global.LbConfig.Port)
	server := http.Server{
		Addr:    addr,
		Handler: Router(),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			global.LbLogger.Fatal(fmt.Sprintf("start http server error: %v", err))
		}
	}()
	logo(addr)

	global.OnExit(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			global.LbLogger.Info(fmt.Sprintf("force exit web：%v", err))
		}
		global.LbLogger.Info("exit web complete！")
		_ = global.LbLogger.Sync() // 确保在程序退出时刷新日志缓冲区
	})

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	global.Exit()
	fmt.Println("Main program stopped")
}

func logo(addr string) {
	fmt.Println("System started, listening: " + addr)
	global.LbLogger.Info("System started, listening: " + addr)
}
