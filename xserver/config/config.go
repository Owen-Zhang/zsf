package config

//Setting 对外提供api接口配制信息
type Setting struct {
	Http Http    `yaml:"http"` //服务启动信息
	Jwt  JwtConf `yaml:"jwt"`  //用户登陆验证信息
}

//Http 服务器信息
type Http struct {
	Mode string `yaml:"mode"` //模式(debug,release,test)
	Port int    `yaml:"port"` //端口
}

//JwtConf 验证用户
type JwtConf struct {
	TimeOut     int64    `yaml:"timeout"`     //超时时间 默认30分钟 1800秒
	Secret      string   `yaml:"secret"`      //jwt签名用到的密钥 默认 qwertyuiop
	ExcludePath []string `yaml:"excludepath"` //不需要验证的路径信息
}

//DefaultConfig api默认配制
func DefaultConfig() *Setting {
	return &Setting{
		Http: Http{
			Mode: "debug",
			Port: 8080,
		},
		Jwt: JwtConf{
			TimeOut:     1800,
			Secret:      "qwertyuiop",
			ExcludePath: []string{},
		},
	}
}
