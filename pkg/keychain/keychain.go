package keychain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"

	"github.com/zalando/go-keyring"
)

const (
	Service = "cloak-cli"
)

func GenerateScopeID(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	hash := sha256.Sum256([]byte(absPath))
	return "project-" + hex.EncodeToString(hash[:8]), nil
}

func Save(scopePath, key string) error {
	user, err := GenerateScopeID(scopePath)
	if err != nil {
		return err
	}
	return keyring.Set(Service, user, key)
}

func Get(scopePath string) (string, error) {
	user, err := GenerateScopeID(scopePath)
	if err != nil {
		return "", err
	}
	return keyring.Get(Service, user)
}

func Delete(scopePath string) error {
	user, err := GenerateScopeID(scopePath)
	if err != nil {
		return err
	}
	return keyring.Delete(Service, user)
}
