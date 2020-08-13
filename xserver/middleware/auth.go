package middleware

import (
	"strings"

	"github.com/Owen-Zhang/zsf/common/xjwt"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/logger"
)

//Auth 验证用户登陆,以及构造用户信息
func auth(secret string, timeOut int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if strings.TrimSpace(token) == "" {
			c.AbortWithStatus(403)
			return
		}
		data, err := xjwt.Decrypt(token, secret)
		if err != nil {
			logger.Errorf("token 解析出错:%+v", err)
			c.AbortWithStatus(403)
			return
		}
		c.Set("login_data", data)
		c.Next()
		tokenNew, err := xjwt.Encrypt(secret, data, timeOut)
		if err != nil {
			logger.Errorf("生成token出错:%+v", err)
			return
		}
		c.Writer.Header().Set("token", tokenNew)
	}
}
