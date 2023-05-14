package pkg

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mall-bff/configs"
)

func NewGrpcConn(ctx context.Context, token string) (*grpc.ClientConn, error) {
	port := configs.Conf.App.SerPort
	addr := fmt.Sprintf("127.0.0.1:%s", port)

	_, ok := ctx.Deadline()
	if !ok {
		// 没有超时的ctx
		return grpc.Dial(addr,
			grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithPerRPCCredentials(&Auth{Token: token}),
		)
	} else {
		// 这里加WithBlock主要想用Context控制连接超时
		return grpc.DialContext(ctx, addr,
			grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithPerRPCCredentials(&Auth{Token: token}),
			grpc.WithBlock(),
		)
	}
}
