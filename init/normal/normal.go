package normal

import (
	"runtime"

	"github.com/Owen-Zhang/zsf/conf"
	"github.com/Owen-Zhang/zsf/logger"
)

//Init 初始系统运行参数
func Init() {
	systemInfo := defaultConfig()
	if err := conf.UnmarshalFile("common.yaml", &systemInfo); err != nil {
		logger.FrameLog.Error(err)
		return
	}
	runtime.GOMAXPROCS(systemInfo.MaxProc)
}
