package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/Owen-Zhang/zsf/config"
	"github.com/toolkits/pkg/logger"
)

//App 应用struct
type App struct {
	log *logger.FileBackend
}

//LoadConfig 加载配制文件
func (a App) LoadConfig() {
	config.Init()
}

//InitLog 初始化日志信息
func (a App) InitLog() {
	lb, err := logger.NewFileBackend("logs")
	if err != nil {
		panic(err)
	}
	a.log = lb
	lb.SetRotateByHour(true)
	lb.SetKeepHours(24)
	logger.SetLogging(logger.ERROR, lb)
}

//initMaxProcs 设置process
func (a App) initMaxProcs() {
	//1 应该去取配制("maxProc")
	maxProc := runtime.NumCPU()

	if maxProc <= 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
		return
	}
	runtime.GOMAXPROCS(maxProc)

}

//clean 清理工作
func (a App) clean() {
	a.log.Flush()
}

func main() {
	app := &App{}
	app.InitLog()
	defer app.clean()

	app.LoadConfig()
	for {
		time.Sleep(10 * time.Second)
		fmt.Println(config.Get("mysql.yaml"))
	}
}
