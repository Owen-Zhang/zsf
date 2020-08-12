package main

import (
	"time"

	"github.com/Owen-Zhang/zsf/config"
	"github.com/Owen-Zhang/zsf/log"
	"github.com/toolkits/pkg/logger"
)

//App 应用struct
type App struct{}

func main() {
	log.Init()
	defer log.Close()

	logger.Error("error from test")
	logger.Warning("warning from test")
	logger.Info("info from test")

	config.Init()
	log.Update()

	go func() {

	}()
	for {
		time.Sleep(10 * time.Second)
		logger.Error("error from test")
		logger.Warning("warning from test")
		logger.Info("info from test")
	}
}
