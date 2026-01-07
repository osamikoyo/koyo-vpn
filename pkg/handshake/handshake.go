package handshake

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
)

var (
	ErrNotSimilar = errors.New("keys is not similar")
)

type HandShakeRouter struct {
	remoteKey string
	selfKey   string

	chiferKey []byte
}

func NewHandShakeRouter(remoteKey, selfKey string) *HandShakeRouter {
	return &HandShakeRouter{
		remoteKey: remoteKey,
		selfKey:   selfKey,
	}
}

func (h *HandShakeRouter) NewHS(REMOTEselfKey string) error {
	if h.remoteKey != REMOTEselfKey {
		return ErrNotSimilar
	}

	key := make([]byte, chacha20poly1305.KeySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return fmt.Errorf("failed generate key: %w", err)
	}

	h.chiferKey = key

	return nil
}

func (h *HandShakeRouter) KeyIsEmpty() bool {
	return len(h.chiferKey) == 0
}

func (h *HandShakeRouter) GetChiferKey() []byte {
	return h.chiferKey
}
