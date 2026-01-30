package repository

import "github.com/pisue/go-playground/grpc-server/config"

type Repository struct {
	cfg *config.Config
}

func NewRepository(cfg *config.Config) (*Repository, error) {
	r := &Repository{cfg: cfg}

	return r, nil
}
