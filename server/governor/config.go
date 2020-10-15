package governor

import (
	"github.com/Owen-Zhang/zsf/util/xnet"
)

//Config 监控配制信息
type Config struct {
	Host   string `yaml:"host"`   //一般为ip地址s
	Port   int    `yaml:"port"`   //端口
	Enable bool   `yaml:"enable"` //是否启用
}

func defaultGovernorConfig() *Config {
	host, _, err := xnet.GetLocalMainIP()
	if err != nil {
		host = "localhost"
	}
	return &Config{
		Host:   host,
		Port:   6681,
		Enable: false,
	}
}
