package log

import (
	"fmt"
	"strings"

	"github.com/Owen-Zhang/zsf/config"
	"github.com/toolkits/pkg/logger"
)

//Init 初始化日志组件
func Init() {
	lb, err := logger.NewFileBackend("logs")
	if err != nil {
		panic(err)
	}
	lb.SetRotateByHour(true)
	lb.SetKeepHours(24)
	logger.SetLogging(logger.ERROR, lb)
}

//Update 更新日志配制信息
//level: "FATAL","ERROR","WARNING","INFO","DEBUG"
func Update() {
	var set setting
	if err := config.UnmarshalFile("logger.yaml", &set); err != nil {
		logger.Error(err)
		return
	}
	if strings.TrimSpace(set.level) != "" {
		fmt.Println(set.level)
		logger.SetSeverity(set.level)
	}
	if set.keepHours > 0 {
		logger.SetKeepHours(set.keepHours)
	}

	logger.SetRotateByHour(set.rotateByHour)
}

//Close 关闭日志(flush到disk)
func Close() {
	logger.Close()
}
