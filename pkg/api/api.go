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
	SenderDeviceID string
	// Payloads maps DeviceIDs to base64 encoded data.
	Payloads map[string]string
}
