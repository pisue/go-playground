package server

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"github.com/pisue/go-playground/grpc-server/config"
	"github.com/pisue/go-playground/grpc-server/gRPC/paseto"
	auth "github.com/pisue/go-playground/grpc-server/gRPC/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	auth.UnimplementedAuthServiceServer
	pasetoMaker    *paseto.PasetoMaker
	tokenVerifyMap map[string]*auth.AuthData
}

func NewGRPCServer(cfg *config.Config) error {
	lis, err := net.Listen("tcp", cfg.GRPC.URL)
	if err != nil {
		return err
	}

	server := grpc.NewServer([]grpc.ServerOption{}...)

	auth.RegisterAuthServiceServer(server, &GRPCServer{
		pasetoMaker:    paseto.NewPasetoMaker(cfg),
		tokenVerifyMap: make(map[string]*auth.AuthData),
	})

	reflection.Register(server)

	go func() {
		log.Println("Start GRPC Sever")
		if err = server.Serve(lis); err != nil {
			panic(err)
		}
	}()

	return nil
}

func (s *GRPCServer) CreateAuth(_ context.Context, req *auth.CreateTokenReq) (*auth.CreateTokenRes, error) {
	data := req.Auth
	token := data.Token

	s.tokenVerifyMap[token] = data

	return &auth.CreateTokenRes{
		Auth: data,
	}, nil
}

func (s *GRPCServer) VerifyAuth(_ context.Context, req *auth.VerifyTokenReq) (*auth.VerifyTokenRes, error) {
	token := req.Token

	res := &auth.VerifyTokenRes{V: &auth.Verify{
		Auth: nil,
	}}

	authData, ok := s.tokenVerifyMap[token]
	if !ok {
		res.V.Status = auth.ResponseType_FAILED
		return res, errors.New("Not Existed At Token VerifyMap")
	} else if err := s.pasetoMaker.VerifyToken(token); err != nil {
		return nil, errors.New("Failed Verify Token")
	} else if authData.ExpireData < time.Now().Unix() {
		delete(s.tokenVerifyMap, token)
		res.V.Status = auth.ResponseType_EXPIRED_DATE
		return res, errors.New("Expired Time Over")
	} else {
		res.V.Status = auth.ResponseType_SUCCESS
		res.V.Auth = authData
		return res, nil
	}

	return res, nil
}
