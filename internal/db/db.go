package db

import "github.com/AYM1607/ccclip/pkg/api"

type DB interface {
	// Users.
	PutUser(id string, passwordHash []byte) error
	GetUser(id string) (*api.User, error)

	// Devices.
	PutDevice(pubKey []byte, userId string) (string, error)
	GetDevice(id string) (*api.Device, error)
	GetUserDevices(userId string) ([]*api.Device, error)
	GetDeviceUser(deviceId string) (*api.User, error)

	// Clipboard.
	PutClipboard(userId string, clipboard *api.Clipboard) error
	GetClipboard(userId string) (*api.Clipboard, error)
}
