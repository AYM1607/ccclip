package crypto

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"golang.org/x/crypto/blake2b"
)

type Direction int

const (
	KeySize uint = 32

	SendDirection Direction = iota
	ReceiveDirection
)

var (
	privateKeyFileMode = os.FileMode(int(0600))
	publicKeyFileMode  = os.FileMode(int(0644))
)

func NewPrivateKey() *ecdh.PrivateKey {
	c := ecdh.X25519()
	k, err := c.GenerateKey(rand.Reader)
	if err != nil {
		panic(fmt.Sprintf("could not generate private key: %s", err.Error()))
	}
	return k
}

func NewSharedKey(local *ecdh.PrivateKey, remote *ecdh.PublicKey, direction Direction) []byte {
	// Calculating a shared secret.
	secret, err := local.ECDH(remote)
	if err != nil {
		panic(err)
	}

	// Take into account public keys for key derivation.
	// See raw_shared_secret @ https://monocypher.org/manual/x25519#DESCRIPTION
	xof, err := blake2b.NewXOF(32, nil)
	if err != nil {
		panic(err)
	}
	xof.Write(secret)
	if direction == SendDirection {
		xof.Write(local.PublicKey().Bytes())
		xof.Write(remote.Bytes())
	} else {
		xof.Write(remote.Bytes())
		xof.Write(local.PublicKey().Bytes())
	}

	key := make([]byte, 32)
	_, err = xof.Read(key)
	if err != nil {
		panic(err)
	}

	return key
}

func PrivateKeyFromB64(encodedKey []byte) *ecdh.PrivateKey {
	keyBytes := make([]byte, KeySize)
	_, err := base64.StdEncoding.Decode(keyBytes, encodedKey)
	if err != nil {
		panic(fmt.Errorf("could not decode base64 private key: %w", err))
	}
	return PrivateKeyFromBytes(keyBytes)
}

func PublicKeyFromB64(encodedKey []byte) *ecdh.PublicKey {
	keyBytes := make([]byte, KeySize)
	_, err := base64.StdEncoding.Decode(keyBytes, encodedKey)
	if err != nil {
		panic(fmt.Errorf("could not decode base64 public key: %w", err))
	}
	return PublicKeyFromBytes(keyBytes)
}

func PrivateKeyFromBytes(keyBytes []byte) *ecdh.PrivateKey {
	key, err := ecdh.X25519().NewPrivateKey(keyBytes)
	if err != nil {
		panic(err)
	}
	return key
}

func PublicKeyFromBytes(keyBytes []byte) *ecdh.PublicKey {
	key, err := ecdh.X25519().NewPublicKey(keyBytes)
	if err != nil {
		panic(err)
	}
	return key
}
