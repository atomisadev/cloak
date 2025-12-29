package store

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/atomisadev/cloak/pkg/crypto"
)

type EncryptedStore map[string]string

func Load(path string, keyHex string) (EncryptedStore, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("store file '%s' not found", path)
	}

	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid key format: %w", err)
	}

	encryptedData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	jsonBytes, err := crypto.Decrypt(encryptedData, key)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	var store EncryptedStore
	if err := json.Unmarshal(jsonBytes, &store); err != nil {
		return nil, fmt.Errorf("corrupted data store: %w", err)
	}

	return store, nil
}

func Save(path string, data EncryptedStore, keyHex string) error {
	key, err := hex.DecodeString(keyHex)
	if err != nil {
		return fmt.Errorf("invalid key format: %w", err)
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	encryptedData, err := crypto.Encrypt(jsonBytes, key)
	if err != nil {
		return err
	}

	return os.WriteFile(path, encryptedData, 0644)
}
