package crypto

import (
	"crypto/rand"

	"golang.org/x/crypto/chacha20poly1305"
)

func Encrypt(key, msg []byte) []byte {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(msg)+aead.Overhead())
	if _, err := rand.Read(nonce); err != nil {
		panic(err)
	}

	return aead.Seal(nonce, nonce, msg, nil)
}

func Decrypt(key, encryptedMsg []byte) []byte {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		panic(err)
	}

	nonce, ciphertext := encryptedMsg[:aead.NonceSize()], encryptedMsg[aead.NonceSize():]
	msg, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err)
	}
	return msg
}
