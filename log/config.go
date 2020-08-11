package log

type setting struct {
	level        string `yaml:"level"`        //记录级别
	rotateByHour bool   `yaml:"rotatebyhour"` //是否一个小时生成一个文件
	keepHours    uint   `yaml:"keephours"`    //日志保存多久时间后删除
}
