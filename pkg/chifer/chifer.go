package chifer

import (
	"crypto/cipher"
	"errors"

	"golang.org/x/crypto/chacha20poly1305"
)

type Chifer struct {
	aead  cipher.AEAD
	nonce []byte 
}


func NewChifer(key string, nonce []byte) (*Chifer, error) {
	aead, err := chacha20poly1305.NewX([]byte(key))
	if err != nil {
		return nil, err
	}

	if len(nonce) != aead.NonceSize() {
		return nil, errors.New("invalid nonce size")
	}

	return &Chifer{
		aead:  aead,
		nonce: append([]byte(nil), nonce...),
	}, nil
}


func (c *Chifer) Encrypt(data []byte) []byte {
	return c.aead.Seal(c.nonce[:0], c.nonce, data, nil)
}


func (c *Chifer) Decrypt(data []byte) ([]byte, error) {
	if len(data) < c.aead.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}



	ciphertextWithTag := data[c.aead.NonceSize():]

	return c.aead.Open(nil, c.nonce, ciphertextWithTag, nil)
}