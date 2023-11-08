package configfile

import (
	"crypto/ecdh"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/AYM1607/ccclip/pkg/crypto"
)

const FileName = "ccclip.yaml"
const PrivateKeyFileName = "private.key"
const PublicKeyFileName = "public.key"

type ConfigFile struct {
	Email    string `yaml:"email"`
	DeviceId string `yaml:"deviceId"`
}

var Path string

func EnsureAndGet() (ConfigFile, error) {
	err := os.MkdirAll(Path, os.FileMode(int(0660)))
	if err != nil {
		return ConfigFile{}, fmt.Errorf("could not create config directory: %w", err)
	}
	rawC, err := os.ReadFile(path.Join(Path, FileName))
	if err != nil {
		if os.IsNotExist(err) {
			return ConfigFile{}, nil
		}
		return ConfigFile{}, fmt.Errorf("could not read current config file: %w", err)
	}
	var c ConfigFile
	return c, json.Unmarshal(rawC, &c)
}

func Write(c ConfigFile) error {
	err := os.MkdirAll(Path, os.FileMode(int(0660)))
	if err != nil {
		return err
	}
	rawC, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("could not convert config to json")
	}
	err = os.WriteFile(path.Join(Path, FileName), rawC, os.FileMode(int(0660)))
	if err != nil {
		return fmt.Errorf("could not write file to config directory: %w", err)
	}
	return nil
}

func LoadPrivateKey() (*ecdh.PrivateKey, error) {
	fp := path.Join(Path, PrivateKeyFileName)
	return crypto.LoadPrivateKey(fp)
}

func LoadPublicKey() (*ecdh.PublicKey, error) {
	fp := path.Join(Path, PublicKeyFileName)
	return crypto.LoadPublicKey(fp)
}

func SavePrivateKey(k *ecdh.PrivateKey) error {
	fp := path.Join(Path, PrivateKeyFileName)
	return crypto.SavePrivateKey(fp, k)
}

func SavePublicKey(k *ecdh.PublicKey) error {
	fp := path.Join(Path, PublicKeyFileName)
	return crypto.SavePublicKey(fp, k)
}
