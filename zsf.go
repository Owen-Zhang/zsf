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
)

//Application 对象实体
type Application struct {
	initOne     sync.Once
	startUpOnce sync.Once
	stopOnce    sync.Once
	serMutex    *sync.RWMutex

	servers []server.IServer
	// api      *xserver.Server
	// governor *governor.Server

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
		app.serMutex = &sync.RWMutex{}
		app.servers = make([]server.IServer, 0)
	})
}

func (app *Application) startUp() {
	app.startUpOnce.Do(func() {
		conf.Init()
		normal.Init()
		log.Init()

	})

	// app := &Application{
	// 	api:      xserver.Init(),
	// 	governor: governor.Init(),
	// }
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

//InitRoute 提供给外部初始化菜单数据
// func (app *Application) InitRoute(route xserver.RouteFunc) {
// 	route(app.api)
// }

//Start 开始运行
func (app *Application) Start() {
	app.api.Start()
	app.waitSignals()
}

//waitSignals 监听系统结束信息
func (app *Application) waitSignals() {
	logger.FrameLog.Info("监听进程结束信号")
	signals.Shutdown(func() {
		app.stop()
	})
}

//stop 退出相关服务
func (app *Application) stop() {
	app.stopOnce.Do(
		func() {
			//关闭api服务
			//关闭注册服务
			//关闭redis连接服务
			//关闭数据库连接服务
			//.....
			app.api.Stop()
			log.Close()
		})
}
