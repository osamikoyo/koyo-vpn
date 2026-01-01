package chifer

import (
	"errors"

	"golang.org/x/crypto/chacha20poly1305"
)

func Encrypt(data, key, nonce []byte) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}

	if len(nonce) != aead.NonceSize() {
		return nil, errors.New("invalid nonce size")
	}

	return aead.Seal(nonce[:0], nonce, data, nil), nil
}

func Decrypt(data, key, nonce []byte) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}

	if len(nonce) != aead.NonceSize() {
		return nil, errors.New("invalid nonce size")
	}

	ciphertextWithTag := data[aead.NonceSize():]

	return aead.Open(nil, nonce, ciphertextWithTag, nil)
}
