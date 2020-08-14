package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/logger"
)

func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Infof("[method:%s];[path:%s];[res-status:%d];[time:%13v]", c.Request.Method, c.Request.URL.Path, c.Writer.Status, time.Now().Sub(start)/time.Second)
	}
}
