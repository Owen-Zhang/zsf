package governor

import (
	"github.com/Owen-Zhang/zsf/conf"
	"github.com/Owen-Zhang/zsf/logger"
)

//Init 加载配制初始化监控信息
func Init() *Server {
	set := defaultGovernorConfig()

	if err := conf.UnmarshalFile("governor.yaml", set); err != nil {
		logger.FrameLog.Errorf("读取监控配制信息出错:%v", err)
		return nil
	}
	if !set.Enable {
		return nil
	}
	return newServer(set)
}
