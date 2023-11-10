package config

type Config struct {
	PublicKeyPath    string
	PrivateKeyPath   string
	DatabaseLocation string
}

var Default = Config{}
