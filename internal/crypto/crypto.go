//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package crypto provides AES-256-GCM encryption for the scratchpad.
//
// The key is a 256-bit random value stored as a raw file. The nonce is
// 12 bytes of random data prepended to the ciphertext. Each write
// re-encrypts the entire file.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
)

// KeySize is the required key length in bytes (256 bits).
const KeySize = 32

// NonceSize is the GCM nonce length in bytes.
const NonceSize = 12

// GenerateKey returns a new 256-bit random key.
//
// Returns:
//   - []byte: A 32-byte random key
//   - error: Non-nil if the system random source fails
func GenerateKey() ([]byte, error) {
	key := make([]byte, KeySize)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, fmt.Errorf("generate key: %w", err)
	}
	return key, nil
}

// Encrypt encrypts plaintext with AES-256-GCM.
//
// The returned ciphertext is formatted as:
//
//	[12-byte nonce][ciphertext + 16-byte GCM tag]
//
// Parameters:
//   - key: 32-byte AES-256 key
//   - plaintext: Data to encrypt
//
// Returns:
//   - []byte: Nonce-prefixed ciphertext
//   - error: Non-nil if key is wrong size or encryption fails
func Encrypt(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %w", err)
	}

	nonce := make([]byte, NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts AES-256-GCM ciphertext produced by [Encrypt].
//
// Parameters:
//   - key: 32-byte AES-256 key (must match the key used for encryption)
//   - ciphertext: Nonce-prefixed ciphertext as produced by [Encrypt]
//
// Returns:
//   - []byte: Decrypted plaintext
//   - error: Non-nil if key is wrong, ciphertext is too short, or
//     authentication fails
func Decrypt(key, ciphertext []byte) ([]byte, error) {
	if len(ciphertext) < NonceSize {
		return nil, errors.New("ciphertext too short")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create GCM: %w", err)
	}

	nonce := ciphertext[:NonceSize]
	data := ciphertext[NonceSize:]

	plaintext, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}

	return plaintext, nil
}

// LoadKey reads a 32-byte key from a file.
//
// Parameters:
//   - path: Path to the key file
//
// Returns:
//   - []byte: The 32-byte key
//   - error: Non-nil if the file cannot be read or is not exactly 32 bytes
func LoadKey(path string) ([]byte, error) {
	key, err := os.ReadFile(path) //nolint:gosec // path is controlled by the caller (config constants)
	if err != nil {
		return nil, fmt.Errorf("read key: %w", err)
	}
	if len(key) != KeySize {
		return nil, fmt.Errorf("invalid key size: got %d bytes, want %d", len(key), KeySize)
	}
	return key, nil
}

// SaveKey writes a key to a file with mode 0600.
//
// Parameters:
//   - path: Destination file path
//   - key: Key bytes to write
//
// Returns:
//   - error: Non-nil if the file cannot be written
func SaveKey(path string, key []byte) error {
	if err := os.WriteFile(path, key, 0600); err != nil {
		return fmt.Errorf("write key: %w", err)
	}
	return nil
}
