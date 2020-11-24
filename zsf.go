package zsf

import (
	"sync"

	"github.com/Owen-Zhang/zsf/conf"
	"github.com/Owen-Zhang/zsf/init/log"
	"github.com/Owen-Zhang/zsf/init/normal"
	"github.com/Owen-Zhang/zsf/logger"
	"github.com/Owen-Zhang/zsf/server"
	"github.com/Owen-Zhang/zsf/server/governor"
	"github.com/Owen-Zhang/zsf/util/signals"
	"github.com/Owen-Zhang/zsf/util/xcycle"
)

//Application 对象实体
type Application struct {
	cycle       *xcycle.Cycle
	initOne     sync.Once
	startUpOnce sync.Once
	stopOnce    sync.Once
	serMutex    *sync.RWMutex

	servers []server.IServer
}

//New 返回实例
func New() *Application {
	app := &Application{}
	app.initialize()
	app.startUp()

	return app
}

func (app *Application) initialize() {
	app.initOne.Do(func() {
		app.cycle = xcycle.NewCycle()
		app.serMutex = &sync.RWMutex{}
		app.servers = make([]server.IServer, 0)
	})
}

func (app *Application) startUp() {
	app.startUpOnce.Do(func() {
		conf.Init()
		normal.Init()
		log.Init()
		app.initGovernor()
	})
}

func (app *Application) initGovernor() error {
	s := governor.Init()
	if s == nil {
		return nil
	}
	app.Serve(s)
	return nil
}

//Serve 加入外部访问服务
func (app *Application) Serve(s ...server.IServer) error {
	app.serMutex.Lock()
	defer app.serMutex.Unlock()
	app.servers = append(app.servers, s...)
	return nil
}

//Run 启动服务,开始运行
func (app *Application) Run() {
	app.waitSignals()
	defer app.clean()
	app.startServers()
	if err := <-app.cycle.Wait(); err != nil {
		logger.FrameLog.Errorf("shutdown with error: %+v", err)
		return
	}
	logger.FrameLog.Info("shutdown normal,bye")
}

//waitSignals 监听系统结束信息
func (app *Application) waitSignals() {
	logger.FrameLog.Info("开始监听进程结束信号")
	signals.Shutdown(func() {
		app.stop()
	})
}

//startServers 启动web服务
func (app *Application) startServers() {
	for _, s := range app.servers {
		s := s
		app.cycle.Run(func() error {
			return s.Start()
		})
	}
}

//clean 结束程序后做一些清理工作
func (app *Application) clean() {
	logger.FrameLog.Close()
	logger.Close()
}

//stop 退出相关服务
func (app *Application) stop() {
	app.stopOnce.Do(
		func() {
			app.serMutex.RLock()
			for _, server := range app.servers {
				server.Stop()
			}
			app.serMutex.RUnlock()
			<-app.cycle.Done()
			app.cycle.Close()
			app.clean()
		})
}
