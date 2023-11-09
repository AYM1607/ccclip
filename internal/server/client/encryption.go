package client

import (
	"crypto/ecdh"

	"github.com/AYM1607/ccclip/pkg/api"
	"github.com/AYM1607/ccclip/pkg/crypto"
)

func encryptForAll(plaintext string, pvk *ecdh.PrivateKey, devices []*api.Device) map[string][]byte {
	res := map[string][]byte{}
	for _, d := range devices {
		key := crypto.NewSharedKey(pvk, crypto.PublicKeyFromBytes(d.PublicKey), crypto.SendDirection)
		res[d.ID] = crypto.Encrypt(key, []byte(plaintext))
	}
	return res
}
