package handshake

import (
	"crypto/rand"
	"errors"
)

var (
	ErrNotSimilar = errors.New("client keys is not similar")
)

type HandShakeRouter struct{
	serverKey string
	clientKey string

	chiferKey string
}

func NewHandShakeRouter(serverKey, clientKey string) *HandShakeRouter {
	return &HandShakeRouter{
		serverKey: serverKey,
		clientKey: clientKey,
	}
}

func (h *HandShakeRouter) NewHS(clientkey string) error {
	if h.clientKey != clientkey {
		return ErrNotSimilar
	}

	chiferkey := rand.Text()
	
	h.chiferKey = chiferkey

	return nil
}

func (h *HandShakeRouter) GetChiferKey() string {
	return h.chiferKey
}