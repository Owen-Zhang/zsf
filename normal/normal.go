package normal

import (
	"runtime"

	"github.com/Owen-Zhang/zsf/config"
	"github.com/Owen-Zhang/zsf/logger"
)

//Init 初始系统运行参数
func Init() {
	systemInfo := defaultConfig()
	if err := config.UnmarshalFile("common.yaml", &systemInfo); err != nil {
		logger.FrameLog.Error(err)
		return
	}
	runtime.GOMAXPROCS(systemInfo.MaxProc)
}
