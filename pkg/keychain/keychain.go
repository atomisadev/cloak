package keychain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/zalando/go-keyring"
)

const (
	Service      = "cloak-cli"
	FallbackDir  = ".cloak"
	FallbackFile = "keystore.json"
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

	if err := keyring.Set(Service, user, key); err == nil {
		return nil
	}

	return saveToLocalStore(user, key)
}

func Get(scopePath string) (string, error) {
	user, err := GenerateScopeID(scopePath)
	if err != nil {
		return "", err
	}

	if key, err := keyring.Get(Service, user); err == nil {
		return key, nil
	}

	return getFromLocalStore(user)
}

func Delete(scopePath string) error {
	user, err := GenerateScopeID(scopePath)
	if err != nil {
		return err
	}

	_ = keyring.Delete(Service, user)
	_ = deleteFromLocalStore(user)

	return nil
}

var storeMutex sync.Mutex

func getStorePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, FallbackDir)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return filepath.Join(dir, FallbackFile), nil
}

func loadStore() (map[string]string, string, error) {
	path, err := getStorePath()
	if err != nil {
		return nil, "", err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return make(map[string]string), path, nil
	}
	if err != nil {
		return nil, path, err
	}

	var store map[string]string
	if err := json.Unmarshal(data, &store); err != nil {
		// If corrupted, return empty to avoid locking user out
		return make(map[string]string), path, nil
	}
	return store, path, nil
}

func saveToLocalStore(user, key string) error {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	store, path, err := loadStore()
	if err != nil {
		return err
	}

	store[user] = key

	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func getFromLocalStore(user string) (string, error) {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	store, _, err := loadStore()
	if err != nil {
		return "", err
	}

	key, ok := store[user]
	if !ok {
		return "", keyring.ErrNotFound
	}
	return key, nil
}

func deleteFromLocalStore(user string) error {
	storeMutex.Lock()
	defer storeMutex.Unlock()

	store, path, err := loadStore()
	if err != nil {
		return err
	}

	delete(store, user)

	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
