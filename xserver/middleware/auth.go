package middleware

import (
	"strings"

	"github.com/Owen-Zhang/zsf/common/xjwt"
	"github.com/Owen-Zhang/zsf/xserver/config"
	"github.com/gin-gonic/gin"
	"github.com/toolkits/pkg/logger"
)

//Auth 验证用户登陆,以及构造用户信息
func Auth(jwtConfig config.JwtConf) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(jwtConfig.ExcludePath) == 0 ||
			isContainPath(jwtConfig.ExcludePath, c.Request.URL.Path) {
			c.Next()
			return
		}
		token := c.GetHeader("token")
		if strings.TrimSpace(token) == "" {
			c.AbortWithStatus(403)
			return
		}
		data, err := xjwt.Decrypt(token, jwtConfig.Secret)
		if err != nil {
			logger.Errorf("token 解析出错:%+v", err)
			c.AbortWithStatus(403)
			return
		}
		c.Set("login_data", data)

		c.Next()

		tokenNew, err := xjwt.Encrypt(jwtConfig.Secret, data, jwtConfig.TimeOut)
		if err != nil {
			logger.Errorf("生成token出错:%+v", err)
			return
		}
		c.Writer.Header().Set("token", tokenNew)
	}
}

//访问的path是否包含在已配制的path中
func isContainPath(confPath []string, inPath string) bool {
	for _, path := range confPath {
		if strings.EqualFold(path, inPath) {
			return true
		}
	}
	return false
}
