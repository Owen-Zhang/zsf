package log

import (
	"github.com/Owen-Zhang/zsf/conf"
	"github.com/Owen-Zhang/zsf/logger"
)

//Init 初始化日志组件
//level: "FATAL","ERROR","WARNING","INFO","DEBUG"
func Init() {
	set := defaultConfig()
	if err := conf.UnmarshalFile("logger.yaml", set); err != nil {
		logger.FrameLog.Errorf("读取日志配制信息出错:%v", err)
		return
	}

	lb, err := logger.NewFileLog("logs/business")
	if err != nil {
		logger.FrameLog.Errorf("new日志记录实例,NewFileLog出错:%v", err)
		return
	}
	lb.SetRotateByHour(set.RotateByHour)
	lb.SetKeepHours(set.KeepHours)
	logger.SetLogging(set.Level, lb)
}

//Close 关闭日志(flush到disk)
func Close() {
	logger.FrameLog.Close()
	logger.Close()
}
