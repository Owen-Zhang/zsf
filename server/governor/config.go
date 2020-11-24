package governor

//Config 监控配制信息
type Config struct {
	Port   int  `yaml:"port"`   //端口
	Enable bool `yaml:"enable"` //是否启用
}

func defaultGovernorConfig() *Config {
	return &Config{
		Port:   6681,
		Enable: false,
	}
}
