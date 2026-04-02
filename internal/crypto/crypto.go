//   /    ctx:                         https://ctx.ist
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
	"io"

	"github.com/ActiveMemory/ctx/internal/config/crypto"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	errCrypto "github.com/ActiveMemory/ctx/internal/err/crypto"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
)

// GenerateKey returns a new 256-bit random key.
//
// Returns:
//   - []byte: A 32-byte random key
//   - error: Non-nil if the system random source fails
func GenerateKey() ([]byte, error) {
	key := make([]byte, crypto.KeySize)
	if _, randErr := io.ReadFull(rand.Reader, key); randErr != nil {
		return nil, errCrypto.GenerateKey(randErr)
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
//   - error: Non-nil if the key is the wrong size or encryption fails
func Encrypt(key, plaintext []byte) ([]byte, error) {
	block, cipherErr := aes.NewCipher(key)
	if cipherErr != nil {
		return nil, errCrypto.CreateCipher(cipherErr)
	}

	gcm, gcmErr := cipher.NewGCM(block)
	if gcmErr != nil {
		return nil, errCrypto.CreateGCM(gcmErr)
	}

	nonce := make([]byte, crypto.NonceSize)
	if _, randErr := io.ReadFull(rand.Reader, nonce); randErr != nil {
		return nil, errCrypto.GenerateNonce(randErr)
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
	if len(ciphertext) < crypto.NonceSize {
		return nil, errCrypto.CiphertextTooShort()
	}

	block, cipherErr := aes.NewCipher(key)
	if cipherErr != nil {
		return nil, errCrypto.CreateCipher(cipherErr)
	}

	gcm, gcmErr := cipher.NewGCM(block)
	if gcmErr != nil {
		return nil, errCrypto.CreateGCM(gcmErr)
	}

	nonce := ciphertext[:crypto.NonceSize]
	data := ciphertext[crypto.NonceSize:]

	plaintext, openErr := gcm.Open(nil, nonce, data, nil)
	if openErr != nil {
		return nil, errCrypto.Decrypt(openErr)
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
	key, readErr := internalIo.SafeReadUserFile(path)
	if readErr != nil {
		return nil, errCrypto.ReadKey(readErr)
	}
	if len(key) != crypto.KeySize {
		return nil, errCrypto.InvalidKeySize(len(key), crypto.KeySize)
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
	if writeErr := internalIo.SafeWriteFile(path, key, fs.PermSecret); writeErr != nil {
		return errCrypto.WriteKey(writeErr)
	}
	return nil
}
