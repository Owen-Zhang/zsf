package main

import (
	"time"

	"github.com/Owen-Zhang/zsf"
	"github.com/Owen-Zhang/zsf/xserver"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/logger"
)

func main() {
	app := zsf.New()
	app.InitRoute(func(s *xserver.Server) {
		user := s.Group("/user")
		user.GET("/info", func(ctx *gin.Context) {
			ctx.Set("result", "aaaaa")
		})
		user.GET("/all", func(ctx *gin.Context) {
			ctx.Set("result", map[string]interface{}{
				"111":  1111,
				"dddd": "ddddd",
			})
		})
	})
	for {
		time.Sleep(10 * time.Second)
		logger.Error("error from test")
		logger.Warning("warning from test")
		logger.Info("info from test")
	}
}
