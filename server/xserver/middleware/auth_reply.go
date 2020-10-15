package middleware

import (
	"github.com/Owen-Zhang/zsf/logger"
	"github.com/Owen-Zhang/zsf/util/xjwt"
	"github.com/Owen-Zhang/zsf/xserver/config"
	"github.com/gin-gonic/gin"
)

//AuthReply 返回数据前返回token
func AuthReply(jwtConfig config.JwtConf) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		data, ok := c.Get("login_data")
		if !ok {
			c.Header(TOKEN, "")
			return
		}
		tokenNew, err := xjwt.Encrypt(jwtConfig.Secret, data, jwtConfig.TimeOut)
		if err != nil {
			logger.FrameLog.Errorf("生成Token出错:%+v", err)
			return
		}
		c.Header(TOKEN, tokenNew)
	}
}
