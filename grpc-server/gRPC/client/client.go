package client

import (
	"github.com/pisue/go-playground/grpc-server/config"
	"github.com/pisue/go-playground/grpc-server/gRPC/paseto"
	auth "github.com/pisue/go-playground/grpc-server/gRPC/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	client      *grpc.ClientConn
	authClient  auth.AuthServiceClient
	pasetoMaker *paseto.PasetoMaker
}

func NewGRPCClient(cfg *config.Config) (*GRPCClient, error) {
	c := new(GRPCClient)

	client, err := grpc.Dial(cfg.GRPC.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	c.client = client
	c.authClient = auth.NewAuthServiceClient(c.client)
	c.pasetoMaker = paseto.NewPasetoMaker(cfg)
	
	return c, nil
}

//rpc CreateAuth(CreateTokenReq) returns (CreateTokenRes);
//rpc VerifyAuth(VerifyTokenReq) returns (VerifyTokenRes);

func (g *GRPCClient) CreateAuth(address string) (*auth.AuthData, error) {
	return nil, nil
}

func (g *GRPCClient) VerifyAuth(token string) (*auth.VerifyTokenRes, error) {
	return nil, nil
}
