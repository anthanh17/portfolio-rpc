package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

func GenerateEncryptionKey(keyLength int) ([]byte, error) {
	key := make([]byte, keyLength)
	if _, err := rand.Read(key); err != nil {
		return nil, fmt.Errorf("error generating random key: %w", err)
	}

	return key, nil
}

func EncryptData(data []byte, key []byte) (string, error) {
	// Create a new AES cipher
	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a GCM cipher
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	// Generate a nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the data
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptData(encodedCiphertext string, key []byte) (string, error) {
	// Decode base64 string
	ciphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		return "", err
	}

	// Create a new AES cipher
	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a GCM cipher
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	// Decrypt the data
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
