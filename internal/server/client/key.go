package client

import (
	"crypto/ecdh"
	"encoding/base64"
	"fmt"

	"github.com/AYM1607/ccclip/pkg/crypto"
)

const serverPublicKeyB64 = "JTyaIVDHe1Nwqmd4NFlkvqj+MZOVp5s3JZP+T3QuoT8="

var serverPublicKey *ecdh.PublicKey

func init() {
	pkeyBytes := make([]byte, crypto.KeySize)
	_, err := base64.StdEncoding.Decode(pkeyBytes, []byte(serverPublicKeyB64))
	if err != nil {
		panic(fmt.Sprintf("cannot decode server public key: %s", err.Error()))
	}
	serverPublicKey = crypto.PublicKeyFromBytes(pkeyBytes)
}
