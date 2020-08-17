package xserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Owen-Zhang/zsf/common/cast"
	cnf "github.com/Owen-Zhang/zsf/config"
	"github.com/Owen-Zhang/zsf/logger"
	"github.com/Owen-Zhang/zsf/xserver/config"
	"github.com/Owen-Zhang/zsf/xserver/middleware"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

//RouteFunc 为注册route用到
type RouteFunc func(*Server)

//Server api服务
type Server struct {
	*gin.Engine                 //gin
	config      *config.Setting //配制信息
	server      *http.Server
}

//Init 对外服务实例化
func Init() *Server {
	set := config.DefaultConfig()
	if err := cnf.UnmarshalFile("server.yaml", set); err != nil {
		logger.FrameLog.Error(err)
	}
	return newserver(set)
}

func newserver(set *config.Setting) *Server {
	ginEngine := gin.New()
	gin.SetMode(set.Http.Mode)
	ginEngine.Use(
		middleware.Recovery(),
		middleware.Log(),
		middleware.Cors(),
		middleware.AuthValid(set.Jwt),
		middleware.Response(),
		middleware.AuthReply(set.Jwt),
	)
	return &Server{
		Engine: ginEngine,
		config: set,
		server: &http.Server{
			ReadTimeout:    60 * time.Second,
			WriteTimeout:   60 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
}

//Start api开始监听服务
func (s *Server) Start() {
	s.server.Addr = fmt.Sprintf(":%d", s.config.Http.Port)
	s.server.Handler = s.Engine

	go func() {
		logger.FrameLog.Infof("开始启动服务,对外端口为: %d", s.config.Http.Port)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.FrameLog.Errorf("服务启动失败:%+v", err)
		}
		logger.FrameLog.Infof("服务启动成功,对外端口为: %d", s.config.Http.Port)
	}()
}

//Stop 直接关闭api服务
func (s *Server) Stop() {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	if err := s.server.Shutdown(ctx); err != nil {
		logger.FrameLog.Fatalf("关闭服务出错:%+v", err)
	}
	select {
	case <-ctx.Done():
		logger.FrameLog.Error("关闭服务等待5秒,时间以过")
	}
	logger.FrameLog.Info("端口服务正常关闭")
}

//SetLoginToken 向客户发送token信息
func SetLoginToken(c *gin.Context, data interface{}) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	tokenBytes, err := json.Marshal(&data)
	if err != nil {
		logger.FrameLog.Errorf("序列化用户登陆信息出错:%+v", err)
		return
	}
	c.Set("login_data", tokenBytes)
}

//GetLoginInfo 获取登陆用户信息
func GetLoginInfo(c *gin.Context, out interface{}) error {
	data, ok := c.Get("login_data")
	if !ok {
		return nil
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Unmarshal([]byte(cast.ToString(data)), out)
}
