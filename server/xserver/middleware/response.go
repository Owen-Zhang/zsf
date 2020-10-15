package middleware

import (
	"net/http"

	"github.com/Owen-Zhang/zsf/util/model"
	"github.com/Owen-Zhang/zsf/util/xerrors"
	"github.com/gin-gonic/gin"
)

//Response 统一处理返回数据
func Response() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if ctx.Writer.Written() {
			return
		}
		data, ok := ctx.Get("result")
		if !ok {
			ctx.AbortWithStatusJSON(
				http.StatusOK,
				model.ResData{
					Status: model.Success,
				})
			return
		}
		switch tmp := data.(type) {
		case string:
			ctx.String(http.StatusOK, tmp)
			return
		case xerrors.BusinessError:
			ctx.JSON(http.StatusOK, model.ResData{
				Status:  model.Fail,
				Message: tmp.Error(),
			})
			return
		case error:
			ctx.JSON(http.StatusOK, model.ResData{
				Status:  model.Fail,
				Message: "访问出错,请联系管理员",
			})
			return

		default:
			ctx.JSON(http.StatusOK, model.ResData{
				Status: model.Success,
				Data:   data,
			})
		}
	}
}
