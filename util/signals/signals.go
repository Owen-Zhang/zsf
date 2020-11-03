package signals

import (
	"os"
	"os/signal"
)

//Shutdown 结束应用
func Shutdown(stop func()) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, shutdownSignals...)

	go func() {
		<-sig
		if stop != nil {
			stop()
		}
		os.Exit(0)
	}()
}
