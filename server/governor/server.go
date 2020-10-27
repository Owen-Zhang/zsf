package governor

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/Owen-Zhang/zsf/logger"
	"github.com/Owen-Zhang/zsf/server"
	"github.com/Owen-Zhang/zsf/util/xnet"
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

//Info 服务信息(此处需要重新构思)
func (s *Server) Info() *server.ServiceInfo {
	ip, _, _ := xnet.GetLocalMainIP()
	return &server.ServiceInfo{
		Name:    "",
		Scheme:  "http",
		Address: ip,
		Port:    8092,
		Enable:  true,
		Group:   "goods",
	}
}

//Start 开启服务
func (s *Server) Start() error {
	err := s.Server.Serve(s.listener)
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	logger.FrameLog.Infof("监控及查看配制服务启动成功,端口为: %d", s.Config.Port)
	return nil
}

//Stop 结束服务
func (s *Server) Stop() error {
	return s.Server.Shutdown(context.Background())
}
