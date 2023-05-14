package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
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

		startSpan := opentracing.GlobalTracer().StartSpan(ctx.Request.URL.Path)
		defer startSpan.Finish()

		ext.HTTPUrl.Set(startSpan, ctx.Request.URL.Path)
		// Http Method
		ext.HTTPMethod.Set(startSpan, ctx.Request.Method)
		// 记录组件名称
		ext.Component.Set(startSpan, "Gin-Http")
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
			startSpan.LogFields(log.String("token", tokens[0]))
			ctx.Set("token", tokens[0])
		}
		ctx.Request = ctx.Request.WithContext(opentracing.ContextWithSpan(ctx.Request.Context(), startSpan))
		ctx.Next()
		ext.HTTPStatusCode.Set(startSpan, uint16(ctx.Writer.Status()))
	}
}
