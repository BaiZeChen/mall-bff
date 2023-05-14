package control

import (
	"context"
	"fmt"
	"github.com/BaiZeChen/mall-api/proto/account"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"google.golang.org/grpc/status"
	"mall-bff/pkg"
	"net/http"
	"time"
)

type Account struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Offset   int    `json:"offset"`
	Limit    int    `json:"limit"`
}

func (a *Account) Login(ctx *gin.Context) {
	err := ctx.ShouldBindJSON(a)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	if len(a.Name) == 0 || len(a.Password) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "参数有误，请重新输入",
		})
		return
	}

	// 链接超时2S
	span, SpanCtx := opentracing.StartSpanFromContext(ctx.Request.Context(), "login")
	defer span.Finish()
	ext.SpanKindRPCClient.Set(span)
	span.LogFields(log.String("name", a.Name), log.String("password", a.Password))

	connCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	conn, err := pkg.NewGrpcConn(connCtx, "")
	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(log.String("err", err.Error()))
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  fmt.Sprintf("请求服务端失败，原因：%s", err.Error()),
		})
		return
	}
	defer conn.Close()
	client := account.NewAccountServiceClient(conn)

	grpcCtx, cancel := context.WithTimeout(SpanCtx, 3*time.Second)
	defer cancel()
	resp, err := client.Login(grpcCtx, &account.ReqAddAccount{
		Name:     a.Name,
		Password: a.Password,
	})
	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(log.String("err", err.Error()))
		fromError, _ := status.FromError(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  fmt.Sprintf("请求服务端失败，原因：%s", fromError.Message()),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "成功",
		"data": gin.H{
			"token": resp.Token,
		},
	})
}

func (a *Account) Add(ctx *gin.Context) {
	err := ctx.ShouldBindJSON(a)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	token, _ := ctx.Get("token")

	if len(a.Name) == 0 || len(a.Name) > 12 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "账户或密码输入有误，请重新输入",
		})
		return
	}
	if len(a.Password) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "账户或密码输入有误，请重新输入",
		})
		return
	}

	// 链接超时2S
	connCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	conn, err := pkg.NewGrpcConn(connCtx, token.(string))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  fmt.Sprintf("请求服务端失败，原因：%s", err.Error()),
		})
		return
	}
	defer conn.Close()
	client := account.NewAccountServiceClient(conn)

	grpcCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err = client.CreateAccount(grpcCtx, &account.ReqAddAccount{
		Name:     a.Name,
		Password: a.Password,
	})
	if err != nil {
		fromError, _ := status.FromError(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  fmt.Sprintf("请求服务端失败，原因：%s", fromError.Message()),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "成功",
	})
}

func (a *Account) UpdateName(ctx *gin.Context) {
	err := ctx.ShouldBindJSON(a)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	token, _ := ctx.Get("token")
	if a.ID == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "没有找到相对应的用户",
		})
		return
	}
	if len(a.Name) == 0 || len(a.Name) > 12 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "账号长度不符合规则，请重新输入！",
		})
		return
	}

	// 链接超时2S
	connCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	conn, err := pkg.NewGrpcConn(connCtx, token.(string))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  fmt.Sprintf("请求服务端失败，原因：%s", err.Error()),
		})
		return
	}
	defer conn.Close()
	client := account.NewAccountServiceClient(conn)

	grpcCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err = client.UpdateAccountName(grpcCtx, &account.ReqUpdateAccountName{
		Id:   uint32(a.ID),
		Name: a.Name,
	})
	if err != nil {
		fromError, _ := status.FromError(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  fmt.Sprintf("请求服务端失败，原因：%s", fromError.Message()),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "成功",
	})

}

func (a *Account) UpdatePassword(ctx *gin.Context) {
	err := ctx.ShouldBindJSON(a)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	token, _ := ctx.Get("token")
	if a.ID == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "没有找到相对应的用户",
		})
		return
	}
	if len(a.Password) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "请填写密码！",
		})
		return
	}

	// 链接超时2S
	connCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	conn, err := pkg.NewGrpcConn(connCtx, token.(string))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  fmt.Sprintf("请求服务端失败，原因：%s", err.Error()),
		})
		return
	}
	defer conn.Close()
	client := account.NewAccountServiceClient(conn)

	grpcCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err = client.UpdateAccountPassword(grpcCtx, &account.ReqUpdateAccountPassword{
		Id:       uint32(a.ID),
		Password: a.Password,
	})
	if err != nil {
		fromError, _ := status.FromError(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  fmt.Sprintf("请求服务端失败，原因：%s", fromError.Message()),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "成功",
	})

}

func (a *Account) Delete(ctx *gin.Context) {
	err := ctx.ShouldBindJSON(a)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	if a.ID == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "没有找到相对应的用户",
		})
		return
	}
	token, _ := ctx.Get("token")

	// 链接超时2S
	connCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	conn, err := pkg.NewGrpcConn(connCtx, token.(string))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  fmt.Sprintf("请求服务端失败，原因：%s", err.Error()),
		})
		return
	}
	defer conn.Close()
	client := account.NewAccountServiceClient(conn)

	grpcCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err = client.DeleteAccount(grpcCtx, &account.ReqDelAccount{Id: uint32(a.ID)})
	if err != nil {
		fromError, _ := status.FromError(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  fmt.Sprintf("请求服务端失败，原因：%s", fromError.Message()),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "成功",
	})
}

func (a *Account) List(ctx *gin.Context) {
	err := ctx.ShouldBindJSON(a)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	if a.Offset < 0 || a.Limit <= 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "页码参数不对",
		})
		return
	}
	token, _ := ctx.Get("token")

	// 链接超时2S
	connCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	conn, err := pkg.NewGrpcConn(connCtx, token.(string))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  fmt.Sprintf("请求服务端失败，原因：%s", err.Error()),
		})
		return
	}
	defer conn.Close()
	client := account.NewAccountServiceClient(conn)

	grpcCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	list, err := client.AccountList(grpcCtx, &account.ReqAccountList{
		Name:   a.Name,
		Offset: uint32(a.Offset),
		Limit:  uint32(a.Limit),
	})
	if err != nil {
		fromError, _ := status.FromError(err)
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  fmt.Sprintf("请求服务端失败，原因：%s", fromError.Message()),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "成功",
		"data": list,
	})
}
