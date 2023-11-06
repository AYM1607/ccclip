package db

import (
	"errors"
	"fmt"

	ulid "github.com/oklog/ulid/v2"

	"github.com/AYM1607/ccclip/pkg/api"
)

type localDB struct {
	users        map[string]*api.User
	devices      map[string]*api.Device
	usersDevices map[string]map[string]struct{}
	clipboards   map[string]*api.Clipboard
}

var _ DB = (*localDB)(nil)

func NewLocalDB() DB {
	return &localDB{
		users:        make(map[string]*api.User),
		devices:      make(map[string]*api.Device),
		usersDevices: make(map[string]map[string]struct{}),
		clipboards:   make(map[string]*api.Clipboard),
	}
}

func (d *localDB) PutUser(id, passwordHash string) error {
	if _, ok := d.users[id]; ok {
		return errors.New("user exists")
	}
	d.users[id] = &api.User{ID: id, PasswordHash: passwordHash}
	return nil
}

func (d *localDB) GetUser(id string) (*api.User, error) {
	u, ok := d.users[id]
	if !ok {
		return nil, errors.New("user does not exist")
	}
	return u, nil
}

func (d *localDB) PutDevice(pubKey []byte, userId string) (string, error) {
	id := ulid.Make().String()
	if _, ok := d.users[userId]; !ok {
		return "", errors.New("user does not exist")
	}

	d.devices[id] = &api.Device{
		ID:        id,
		PublicKey: pubKey,
	}

	if d.usersDevices[userId] == nil {
		d.usersDevices[userId] = map[string]struct{}{}
	}
	d.usersDevices[userId][id] = struct{}{}
	return id, nil
}

func (d *localDB) GetDevice(id string) (*api.Device, error) {
	device, ok := d.devices[id]
	if !ok {
		return nil, errors.New("requested device does no exist")
	}
	return device, nil
}

func (d *localDB) GetUserDevices(userId string) ([]*api.Device, error) {
	ids := d.usersDevices[userId]
	res := []*api.Device{}
	for id := range ids {
		d, ok := d.devices[id]
		if !ok {
			return nil, fmt.Errorf("device %s is associated to user but it does not exist", id)
		}
		res = append(res, d)
	}
	return res, nil
}

func (d *localDB) GetDeviceUser(deviceId string) (*api.User, error) {
	// Foreign keys can make this better.
	for uId, u := range d.users {
		if _, ok := d.usersDevices[uId]; ok {
			return u, nil
		}
	}
	return nil, errors.New("device is not associated to any user")
}

func (d *localDB) PutClipboard(userId string, clipboard *api.Clipboard) error {
	d.clipboards[userId] = clipboard
	return nil
}

func (d *localDB) GetClipboard(userId string) (*api.Clipboard, error) {
	c, ok := d.clipboards[userId]
	if !ok {
		return nil, errors.New("user does not have a current clipboard")
	}
	return c, nil
}
