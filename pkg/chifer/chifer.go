package chifer

import (
	"encoding/binary"
	"errors"

	"golang.org/x/crypto/chacha20poly1305"
)

func Encrypt(data, key []byte, nonceUint64 uint64) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}

	var nonce [12]byte
	binary.LittleEndian.PutUint64(nonce[0:8], nonceUint64)

	return aead.Seal(nil, nonce[:], data, nil), nil
}

func Decrypt(ciphertext, key []byte) ([]byte, error) {
	if len(ciphertext) < chacha20poly1305.NonceSize+chacha20poly1305.Overhead {
		return nil, errors.New("ciphertext too short")
	}

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}

	nonce := ciphertext[:chacha20poly1305.NonceSize]
	ctWithTag := ciphertext[chacha20poly1305.NonceSize:]

	return aead.Open(nil, nonce, ctWithTag, nil)
}