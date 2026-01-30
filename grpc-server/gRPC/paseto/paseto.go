package paseto

import (
	"crypto/rand"

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
	randomBytes := make([]byte, 16)
	r, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return m.Pt.Encrypt(m.Key, auth, r)
}

func (m *PasetoMaker) VerifyToken(token string) error {
	var a *auth.AuthData
	return m.Pt.Decrypt(token, m.Key, a, nil)
}
