package middleware

import (
	"time"

	"github.com/Owen-Zhang/zsf/logger"
	"github.com/gin-gonic/gin"
)

//Log 记录访问的地址、时间等
func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.FrameLog.Infof("[method:%s];[path:%s];[status:%d];[time:%d (ms)]", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), time.Since(start)/time.Millisecond)
	}
}
