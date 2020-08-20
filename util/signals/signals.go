package signals

import (
	"fmt"
	"os"
	"os/signal"
)

//Shutdown 结束应用
func Shutdown(stop func()) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, shutdownSignals...)

	select {
	case <-sig:
		fmt.Printf("捕获到结束程序信号, 正在线束中... pid=%d\n", os.Getpid())
	}
	if stop != nil {
		stop()
	}
	os.Exit(0)
}
