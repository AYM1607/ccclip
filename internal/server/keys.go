package server

import (
	"crypto/ecdh"
	"errors"
	"os"

	"github.com/AYM1607/ccclip/pkg/crypto"
)

const (
	privateKeyEnv     = "CCCLIP_PRIVATE_KEY"
	publicKeyEnv      = "CCCLIP_PUBLIC_KEY"
	privateKeyPathEnv = "CCCLIP_PRIVATE_KEY_PATH"
	publicKeyPathEnv  = "CCCLIP_PUBLIC_KEY_PATH"
)

func loadKeys() (*ecdh.PrivateKey, *ecdh.PublicKey, error) {
	// Prioritize explicit keys over files.
	var pvk *ecdh.PrivateKey
	var pbk *ecdh.PublicKey

	if b64PrivateKey := os.Getenv(privateKeyEnv); b64PrivateKey != "" {
		pvk = crypto.PrivateKeyFromB64([]byte(b64PrivateKey))
	} else if privateKeyPath := os.Getenv(privateKeyPathEnv); privateKeyPath != "" {
		pvk = crypto.LoadPrivateKeyFromFile(privateKeyPath)
	} else {
		return nil, nil, errors.New("no private key was found")
	}

	if b64PublicKey := os.Getenv(publicKeyEnv); b64PublicKey != "" {
		pbk = crypto.PublicKeyFromB64([]byte(b64PublicKey))
	} else if publicKeyPath := os.Getenv(publicKeyPathEnv); publicKeyPath != "" {
		pbk = crypto.LoadPublicKeyFromFile(publicKeyPath)
	} else {
		return nil, nil, errors.New("to public key was found")
	}

	return pvk, pbk, nil
}
