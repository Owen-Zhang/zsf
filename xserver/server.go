package xserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Owen-Zhang/zsf/common/cast"
	"github.com/Owen-Zhang/zsf/xserver/middleware"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/toolkits/pkg/logger"
)

//RouteFunc 为注册route用到
type RouteFunc func(*Server)

//Server api服务
type Server struct {
	*gin.Engine          //gin
	config      *Setting //配制信息
	server      *http.Server
}

func newserver(set *Setting) *Server {
	return &Server{
		Engine: gin.New(),
		config: set,
		server: &http.Server{
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
}

//Start api开始监听服务
func (s *Server) Start() {
	gin.SetMode(s.config.Http.Mode)
	s.Engine.Use(
		middleware.Recovery(),
		middleware.Log(),
		middleware.Cors(),
		middleware.Auth(s.config.Jwt),
		middleware.Response(),
	)
	s.server.Addr = fmt.Sprintf(":%d", s.config.Http.Port)
	s.server.Handler = s.Engine

	go func() {
		logger.Infof("开始启动服务,对外端口为: %d", s.config.Http.Port)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("服务启动失败:%+v", err)
		}
	}()
}

//Shutdown 关闭端口服务
func (s *Server) Shutdown() {
	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	if err := s.server.Shutdown(ctx); err != nil {
		logger.Fatalf("关闭服务出错:%+v", err)
	}
	select {
	case <-ctx.Done():
		logger.Error("关闭服务等待5秒,时间以过")
	}
	logger.Info("端口服务正常关闭")
}

//SetLoginToken 向客户发送token信息
func SetLoginToken(c *gin.Context, data interface{}) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	tokenBytes, err := json.Marshal(&data)
	if err != nil {
		logger.Errorf("序列化用户登陆信息出错:%+v", err)
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
