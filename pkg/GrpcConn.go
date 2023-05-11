package pkg

import (
	"github.com/BaiZeChen/mall-api/proto/account"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mall-bff/configs"
)

type GrpcConn struct {
	Client account.AccountServiceClient
	Conn   *grpc.ClientConn
}

func (g *GrpcConn) NewConn(token string) error {
	port := configs.Conf.App.SerPort
	conn, err := grpc.Dial("127.0.0.1:"+port, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithPerRPCCredentials(&Auth{Token: token}))
	if err != nil {
		return err
	}

	g.Conn = conn
	g.Client = account.NewAccountServiceClient(conn)
	return nil
}
