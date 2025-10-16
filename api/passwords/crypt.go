package passwords

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func BCrypt(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return "{BCRYPT}" + string(hashed), nil
}

func VerifyBCrypt(password, hashed string) bool {
	hashed = strings.TrimPrefix(hashed, "{BCRYPT}")
	if strings.HasPrefix(hashed, "{CRYPT}") {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}

func Encrypt(text string) (string, error) {
	key := os.Getenv("CRYPT_KEY")
	if len(key) < 32 {
		return "", fmt.Errorf("invalid key length: must be greater than or equal to 32 bytes but got %d", len(key))
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nil, nonce, []byte(text), nil)
	result := append(nonce, ciphertext...)
	encoded := base64.StdEncoding.EncodeToString(result)

	return "{AES}" + encoded, nil
}

func Decrypt(encrypted string) (string, error) {
	encrypted = strings.TrimPrefix(encrypted, "{AES}")
	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode: %w", err)
	}

	if len(data) < 12 {
		return "", fmt.Errorf("invalid encrypted data: too short")
	}

	key := os.Getenv("CRYPT_KEY")
	if len(key) < 32 {
		return "", fmt.Errorf("invalid key length: must be greater than or equal to 32 bytes but got %d", len(key))
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce, ciphertext := data[:12], data[12:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}
