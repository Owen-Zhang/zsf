package redis

//Config redis配制信息
type Config struct {
	// Addrs 实例配置地址
	Addrs []string `yaml:"addrs"`
	// Password 密码
	Password string `yaml:"password"`
	// ReadTimeout 读超时 默认3s
	ReadTimeout int `yaml:"readtimeout"`
	// WriteTimeout 读超时 默认3s
	WriteTimeout int `yaml:"writetimeout"`
}

func defaultRedisConfig() *Config {
	return &Config{
		ReadTimeout:  3,
		WriteTimeout: 3,
	}
}
