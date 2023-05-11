package control

import (
	"context"
	"fmt"
	"github.com/BaiZeChen/mall-api/proto/account"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
	"mall-bff/pkg"
	"net/http"
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

	grpcClient := &pkg.GrpcConn{}
	err = grpcClient.NewConn("")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	defer grpcClient.Conn.Close()
	resp, err := grpcClient.Client.Login(context.Background(), &account.ReqAddAccount{
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

	grpcClient := &pkg.GrpcConn{}
	err = grpcClient.NewConn(token.(string))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	defer grpcClient.Conn.Close()
	_, err = grpcClient.Client.CreateAccount(context.Background(), &account.ReqAddAccount{
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

	grpcClient := &pkg.GrpcConn{}
	err = grpcClient.NewConn(token.(string))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	defer grpcClient.Conn.Close()
	_, err = grpcClient.Client.UpdateAccountName(context.Background(), &account.ReqUpdateAccountName{
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

	grpcClient := &pkg.GrpcConn{}
	err = grpcClient.NewConn(token.(string))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	defer grpcClient.Conn.Close()
	_, err = grpcClient.Client.UpdateAccountPassword(context.Background(), &account.ReqUpdateAccountPassword{
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

	grpcClient := &pkg.GrpcConn{}
	err = grpcClient.NewConn(token.(string))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	defer grpcClient.Conn.Close()

	_, err = grpcClient.Client.DeleteAccount(context.Background(), &account.ReqDelAccount{Id: uint32(a.ID)})
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

	grpcClient := &pkg.GrpcConn{}
	err = grpcClient.NewConn(token.(string))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  err.Error(),
		})
		return
	}
	defer grpcClient.Conn.Close()

	list, err := grpcClient.Client.AccountList(context.Background(), &account.ReqAccountList{
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
