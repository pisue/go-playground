package network

import (
	"github.com/gin-gonic/gin"
	"github.com/pisue/go-playground/grpc-server/config"
	"github.com/pisue/go-playground/grpc-server/gRPC/client"
	"github.com/pisue/go-playground/grpc-server/service"
)

type Network struct {
	cfg *config.Config

	service    *service.Service
	gRPCClient *client.GRPCClient

	engin *gin.Engine
}

func NewNetwork(cfg *config.Config, service *service.Service, gRPCClient *client.GRPCClient) (*Network, error) {
	r := &Network{cfg: cfg, service: service, engin: gin.New(), gRPCClient: gRPCClient}

	// 1. token 생성하는 API
	r.engin.POST("/login", r.login)

	// 2. token 검증하는 API
	r.engin.GET("/verify", r.verifyLogin(), r.verify)

	return r, nil
}

func (n *Network) StartServer() {
	n.engin.Run(":8080")
}
