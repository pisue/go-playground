package paseto

import (
	"github.com/o1egl/paseto"
	"github.com/pisue/go-playground/grpc-server/config"
	auth "github.com/pisue/go-playground/grpc-server/gRPC/proto"
)

type PasetoMaker struct {
	Pt  *paseto.V2
	Key []byte
}

func NewPasetoMaker(cfg *config.Config) *PasetoMaker {
	return &PasetoMaker{
		Pt:  paseto.NewV2(),
		Key: []byte(cfg.Pasteo.Key),
	}
}

func (m *PasetoMaker) CreateNewToken(auth *auth.AuthData) (string, error) {
	return "", nil
}

func (m *PasetoMaker) VerifyToken(token string) error {
	return nil
}
