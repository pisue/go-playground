package client

import (
	"context"
	"time"

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

func (g *GRPCClient) CreateAuth(name string) (*auth.AuthData, error) {
	now := time.Now()
	expiredTime := now.Add(30 * time.Minute)

	a := &auth.AuthData{
		Name:       name,
		CreateData: now.Unix(),
		ExpireData: expiredTime.Unix(),
	}

	token, err := g.pasetoMaker.CreateNewToken(a)
	if err != nil {
		return nil, err
	}

	a.Token = token

	res, err := g.authClient.CreateAuth(context.Background(), &auth.CreateTokenReq{Auth: a})
	if err != nil {
		return nil, err
	}

	return res.Auth, nil
}

func (g *GRPCClient) VerifyAuth(token string) (*auth.Verify, error) {
	res, err := g.authClient.VerifyAuth(context.Background(), &auth.VerifyTokenReq{Token: token})
	if err != nil {
		return nil, err
	}

	return res.V, nil
}
