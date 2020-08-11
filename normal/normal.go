package normal

import (
	"runtime"

	"github.com/Owen-Zhang/zsf/config"
	"github.com/toolkits/pkg/logger"
)

//Init 初始系统运行参数
func Init() {
	var systemInfo systemP
	if err := config.UnmarshalFile("common", &systemInfo); err != nil {
		logger.Error(err)
		return
	}

	maxProc := systemInfo.maxProc
	if maxProc <= 0 {
		maxProc = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(maxProc)
}
