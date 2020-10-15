package governor

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/Owen-Zhang/zsf/logger"
)

var (
	DefaultServeMux = http.NewServeMux()
	routes          = []string{}
)

//Server Server
type Server struct {
	*http.Server
	listener net.Listener
	*Config
}

func newServer(config *Config) *Server {
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)
	lister, err := net.Listen("tcp4", address)
	if err != nil {
		logger.Fatalf("启动governor时出现错误:%v", err)
	}

	return &Server{
		Server: &http.Server{
			Addr:    address,
			Handler: DefaultServeMux,
		},
		listener: lister,
		Config:   config,
	}
}

//Start 开启服务
func (s *Server) Start() error {
	err := s.Server.Serve(s.listener)
	if err != nil && err != http.ErrServerClosed {
		return err
		//logger.FrameLog.Fatalf("服务启动失败:%+v", err)
	}
	logger.FrameLog.Infof("监控及查看配制服务启动成功,端口为: %d", s.Config.Port)
	return nil
}

//Stop 结束服务
func (s *Server) Stop() error {
	return s.Server.Shutdown(context.Background())
}
