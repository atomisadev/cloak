package sharing

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/atomisadev/cloak/pkg/crypto"
)

const FileIOEndpoint = "https://file.io"

type FileIOResponse struct {
	Success bool   `json:"success"`
	Key     string `json:"key"`
	Link    string `json:"link"`
	Message string `json:"message,omitempty"`
}

func CreateDeadDrop(masterKey []byte) (string, error) {
	fragmentKeyHex, err := crypto.GenerateKey()
	if err != nil {
		return "", fmt.Errorf("failed to generate fragment key: %w", err)
	}
	fragmentKey, _ := hex.DecodeString(fragmentKeyHex)

	encryptedBlob, err := crypto.Encrypt(masterKey, fragmentKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt master key: %w", err)
	}

	fileKey, err := uploadToEphemeralStore(encryptedBlob)
	if err != nil {
		return "", fmt.Errorf("upload failed: %w", err)
	}

	magicLink := fmt.Sprintf("https://cloak.hitmo.xyz/claim/%s#%s", fileKey, fragmentKeyHex)

	return magicLink, nil
}

func uploadToEphemeralStore(data []byte) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "secret.bin")
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, bytes.NewReader(data)); err != nil {
		return "", err
	}

	_ = writer.WriteField("expires", "1d")

	if err := writer.Close(); err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", FileIOEndpoint, body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	var result FileIOResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if !result.Success {
		return "", fmt.Errorf("api error: %s", result.Message)
	}

	return result.Key, nil
}
