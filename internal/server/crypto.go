package server

import (
	"crypto/ecdh"
	"encoding/json"
	"time"

	"github.com/AYM1607/ccclip/internal/db"
	"github.com/AYM1607/ccclip/pkg/crypto"
)

type AuthenticatedPayload struct {
	DeviceID string `json:"deviceID"`
	Payload  []byte `json:"payload"`
}

type FingerPrint struct {
	Timestamp time.Time `json:"timestamp"`
}

func decryptAuthenticatedPayload[T any](p AuthenticatedPayload, d db.DB, pk *ecdh.PrivateKey) (T, error) {
	var res T
	var zero T

	device, err := d.GetDevice(p.DeviceID)
	if err != nil {
		return zero, err
	}

	key := crypto.NewSharedKey(pk, crypto.PublicKeyFromBytes(device.PublicKey), crypto.ReceiveDirection)
	plain := crypto.Decrypt(key, p.Payload)

	err = json.Unmarshal(plain, res)
	if err != nil {
		return zero, err
	}
	return res, nil
}
