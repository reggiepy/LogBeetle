package boot

import (
	"context"
	"errors"
	"fmt"
	"github.com/reggiepy/LogBeetle/global"
	"github.com/reggiepy/goutils/signailUtils"
	"net/http"
	"time"
)

func Boot() {
	addr := fmt.Sprintf("%s:%d", "0.0.0.0", global.LbConfig.Port)
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

	signailUtils.OnExit(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			global.LbLogger.Info(fmt.Sprintf("force exit web：%v", err))
		}
		global.LbLogger.Info("exit web complete！")
		global.LbLoggerClearup() // 确保在程序退出时刷新日志缓冲区
	})

	signailUtils.WaitExit(1 * time.Second)
}

func logo(addr string) {
	global.LbLogger.Info("System started, listening: http://" + addr)
}
