package boot

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/reggiepy/LogBeetle/global"
	"go.uber.org/zap"
)

func Boot() {
	addr := fmt.Sprintf("%s:%d", "127.0.0.1", global.LbConfig.Port)
	server := http.Server{
		Addr:    addr,
		Handler: Router(),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			global.LbLogger.Info(fmt.Sprintf("start http server error: %v\n", err))
		}
	}()
	logo(addr)
}

func logo(addr string) {
	fmt.Println("System started, listening: " + addr)
	zap.L().Info("System started, listening: " + addr)
}
