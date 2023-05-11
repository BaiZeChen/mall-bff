package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		arr := strings.Split(path, "/")
		if len(arr) == 0 {
			ctx.Abort()
			ctx.JSON(http.StatusNotFound, gin.H{
				"code": 0,
				"msg":  "请求不存在",
			})
			return
		}
		path = arr[len(arr)-1]
		if path == "" {
			ctx.Abort()
			ctx.JSON(http.StatusNotFound, gin.H{
				"code": 0,
				"msg":  "请求不存在",
			})
			return
		}
		if path != "login" {
			// 登录不用验证token
			tokens, ok := ctx.Request.Header["Token"]
			if !ok {
				ctx.Abort()
				ctx.JSON(http.StatusOK, gin.H{
					"code": 0,
					"msg":  "请登录",
				})
				return
			}
			ctx.Set("token", tokens[0])
		}
		ctx.Next()
	}
}
