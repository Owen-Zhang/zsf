package xserver

import (
	"github.com/Owen-Zhang/zsf/config"
	"github.com/toolkits/pkg/logger"
)

//Setting 对外提供api接口配制信息
type Setting struct {
	Http Http `yaml:"http"` //服务启动信息
	Jwt  Jwt  `yaml:"jwt"`  //用户登陆验证信息
}

//Http 服务器信息
type Http struct {
	Mode string `yaml:"mode"` //模式(debug,release,test)
	Port int    `yaml:"port"` //端口
}

//Jwt 验证用户
type Jwt struct {
	TimeOut int64  `yaml:"timeout"` //超时时间 默认30分钟 1800秒
	Secret  string `yaml:"secret"`  //jwt签名用到的密钥 默认 qwertyuiop
}

//Init 对外服务实例化
func Init() {
	set := defaultConfig()
	if err := config.UnmarshalFile("server.yaml", set); err != nil {
		logger.Error(err)
	}
}

func defaultConfig() *Setting {
	return &Setting{
		Http: Http{
			Mode: "debug",
			Port: 8080,
		},
		Jwt: Jwt{
			TimeOut: 1800,
			Secret:  "qwertyuiop",
		},
	}
}
