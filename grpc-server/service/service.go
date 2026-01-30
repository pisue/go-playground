package service

import (
	"github.com/pisue/go-playground/grpc-server/config"
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
