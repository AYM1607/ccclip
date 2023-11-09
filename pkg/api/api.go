package api

type User struct {
	ID           string
	PasswordHash string
}

type Device struct {
	ID        string
	PublicKey []byte
}

type Clipboard struct {
	SenderPublicKey []byte
	// Payloads maps DeviceIDs to encrypted data.
	Payloads map[string][]byte
}
