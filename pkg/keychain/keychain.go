package keychain

import "github.com/zalando/go-keyring"

const (
	Service = "cloak-cli"
	User    = "master-key"
)

func Save(key string) error {
	return keyring.Set(Service, User, key)
}

func Get() (string, error) {
	return keyring.Get(Service, User)
}

func Delete() error {
	return keyring.Delete(Service, User)
}
