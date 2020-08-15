package zsf

import (
	"github.com/Owen-Zhang/zsf/config"
	"github.com/Owen-Zhang/zsf/log"
	"github.com/Owen-Zhang/zsf/xserver"
)

//Application 对象实体
type Application struct {
	api *xserver.Server
}

//New 返回实例
func New() *Application {
	config.Init()
	log.Init()
	app := &Application{
		api: xserver.Init(),
	}
	return app
}

//InitRoute 提供给外部初始化菜单数据
func (app *Application) InitRoute(route xserver.RouteFunc) {
	route(app.api)
}

//Start 开始运行
func (app *Application) Start() {
	app.api.Start()
}

//Close 关闭要处理的一些事情
func (app *Application) Close() {
	log.Close()
}
