package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"mall-bff/configs"
	"mall-bff/internal/control"
	"mall-bff/internal/middleware"
	"mall-bff/pkg"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	pkg.InitTracing()

	engine := gin.Default()
	engine.Use(middleware.Auth())

	RegisterRoute(engine)

	app := configs.Conf.App
	addr := fmt.Sprintf(":%s", app.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: engine,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("服务器启动失败，原因：%s\n", err))
		}
	}()
	fmt.Println("服务器已启动")

	Shutdown(srv)
}

func Shutdown(ser *http.Server) {
	sign := make(chan os.Signal)
	signal.Notify(sign, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	select {
	case <-sign:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := ser.Shutdown(ctx); err != nil {
			panic(err)
		}
		err := pkg.TracingCloser.Close()
		if err != nil {
			panic(err)
		}
		fmt.Println("服务已关闭！")
	}
}

func RegisterRoute(engine *gin.Engine) {

	group := engine.Group("/mall/bff")
	{
		userGroup := group.Group("/account")
		{
			userGroup.POST("/login", func(c *gin.Context) {
				accountControl := &control.Account{}
				accountControl.Login(c)
			})
			userGroup.POST("/add", func(c *gin.Context) {
				accountControl := &control.Account{}
				accountControl.Add(c)
			})
			userGroup.POST("/update/name", func(c *gin.Context) {
				accountControl := &control.Account{}
				accountControl.UpdateName(c)
			})
			userGroup.POST("/update/password", func(c *gin.Context) {
				accountControl := &control.Account{}
				accountControl.UpdatePassword(c)
			})
			userGroup.POST("/delete", func(c *gin.Context) {
				accountControl := &control.Account{}
				accountControl.Delete(c)
			})
			userGroup.POST("/list", func(c *gin.Context) {
				accountControl := &control.Account{}
				accountControl.List(c)
			})
		}
	}
}
