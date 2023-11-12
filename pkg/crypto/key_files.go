package crypto

import (
	"crypto/ecdh"
	"encoding/base64"
	"os"
)

func LoadPrivateKeyFromFile(fp string) *ecdh.PrivateKey {
	kb, err := loadKeyFile(fp)
	if err != nil {
		panic(err)
	}
	return PrivateKeyFromBytes(kb)
}

func LoadPublicKeyFromFile(fp string) *ecdh.PublicKey {
	kb, err := loadKeyFile(fp)
	if err != nil {
		panic(err)
	}

	return PublicKeyFromBytes(kb)
}

func SavePrivateKeyToFile(fp string, k *ecdh.PrivateKey) error {
	return saveKeyFile(fp, k.Bytes(), privateKeyFileMode)
}

func SavePublicKeyToFile(fp string, k *ecdh.PublicKey) error {
	return saveKeyFile(fp, k.Bytes(), publicKeyFileMode)
}

func loadKeyFile(fp string) ([]byte, error) {
	b64Key, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	keyBytes := make([]byte, KeySize)
	base64.StdEncoding.Decode(keyBytes, b64Key)
	return keyBytes, nil
}

func saveKeyFile(fp string, key []byte, fm os.FileMode) error {
	b64Key := make([]byte, base64.StdEncoding.EncodedLen(len(key)))
	base64.StdEncoding.Encode(b64Key, key)

	return os.WriteFile(fp, b64Key, fm)
}
