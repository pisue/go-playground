package service

import (
	"github.com/pisue/go-playground/grpc-server/config"
	auth "github.com/pisue/go-playground/grpc-server/gRPC/proto"
	"github.com/pisue/go-playground/grpc-server/repository"
)

type Service struct {
	cfg        *config.Config
	repository *repository.Repository
}

func NewService(cfg *config.Config, repository *repository.Repository) (*Service, error) {
	r := &Service{cfg: cfg, repository: repository}

	return r, nil
}

func (s *Service) CreateAuth(name string) (*auth.AuthData, error) {
	return s.repository.CreateAuth(name)
}
