package log

//Setting level: "FATAL","ERROR","WARNING","INFO","DEBUG"
type Setting struct {
	Level        string `yaml:"level"`        //记录级别
	RotateByHour bool   `yaml:"rotatebyhour"` //是否一个小时生成一个文件
	KeepHours    uint   `yaml:"keephours"`    //日志保存多久时间后删除
}

func defaultConfig() *Setting {
	return &Setting{
		Level:        "ERROR",
		RotateByHour: true,
		KeepHours:    24,
	}
}
