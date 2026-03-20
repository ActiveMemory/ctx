//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package crypto

import (
	"errors"
	"fmt"
	"os"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// LoadKey classifies a key-loading failure.
//
// If the underlying error is os.ErrNotExist, returns NoKeyAt(keyPath).
// Otherwise wraps the cause as a generic load-key error.
//
// Parameters:
//   - cause: the underlying error from crypto.LoadKey
//   - keyPath: the resolved key path that was checked
//
// Returns:
//   - error: NoKeyAt or "load key: <cause>"
func LoadKey(cause error, keyPath string) error {
	if errors.Is(cause, os.ErrNotExist) {
		return NoKeyAt(keyPath)
	}
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrCryptoLoadKey), cause,
	)
}

// EncryptFailed wraps an encryption failure.
//
// Parameters:
//   - cause: the underlying error from crypto.Encrypt.
//
// Returns:
//   - error: "encrypt: <cause>"
func EncryptFailed(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrCryptoEncryptFailed), cause,
	)
}

// DecryptFailed returns an error indicating decryption failure.
//
// Returns:
//   - error: "decryption failed: wrong key?"
func DecryptFailed() error {
	return errors.New(desc.TextDesc(text.DescKeyErrCryptoDecryptFailed))
}

// NoKeyAt returns an error indicating a missing encryption key.
//
// Parameters:
//   - path: the resolved key path that was checked.
//
// Returns:
//   - error: "encrypted scratchpad found but no key at <path>"
func NoKeyAt(path string) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrCryptoNoKeyAt), path,
	)
}

// SaveKey wraps a failure to save an encryption key.
//
// Parameters:
//   - cause: the underlying error from key saving
//
// Returns:
//   - error: "failed to save scratchpad key: <cause>"
func SaveKey(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrCryptoSaveKey), cause)
}

// MkdirKeyDir wraps a failure to create the key directory.
//
// Parameters:
//   - cause: the underlying OS error
//
// Returns:
//   - error: "failed to create key dir: <cause>"
func MkdirKeyDir(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrCryptoMkdirKeyDir), cause,
	)
}

// CreateCipher wraps a failure to create an AES cipher.
//
// Parameters:
//   - cause: the underlying crypto error.
//
// Returns:
//   - error: "create cipher: <cause>"
func CreateCipher(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrCryptoCreateCipher), cause,
	)
}

// CreateGCM wraps a failure to create a GCM instance.
//
// Parameters:
//   - cause: the underlying crypto error.
//
// Returns:
//   - error: "create GCM: <cause>"
func CreateGCM(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrCryptoCreateGCM), cause,
	)
}

// GenerateNonce wraps a failure to generate a random nonce.
//
// Parameters:
//   - cause: the underlying IO error.
//
// Returns:
//   - error: "generate nonce: <cause>"
func GenerateNonce(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrCryptoGenerateNonce), cause,
	)
}

// GenerateKey wraps a failure to generate a random key.
//
// Parameters:
//   - cause: the underlying IO error.
//
// Returns:
//   - error: "generate key: <cause>"
func GenerateKey(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrCryptoGenerateKey), cause,
	)
}

// CiphertextTooShort returns an error when ciphertext is shorter
// than the nonce size.
//
// Returns:
//   - error: "ciphertext too short"
func CiphertextTooShort() error {
	return errors.New(
		desc.TextDesc(text.DescKeyErrCryptoCiphertextTooShort),
	)
}

// Decrypt wraps a decryption failure with cause.
//
// Parameters:
//   - cause: the underlying decryption error.
//
// Returns:
//   - error: "decrypt: <cause>"
func Decrypt(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrCryptoDecrypt), cause,
	)
}

// ReadKey wraps a failure to read a key file.
//
// Parameters:
//   - cause: the underlying read error.
//
// Returns:
//   - error: "read key: <cause>"
func ReadKey(cause error) error {
	return fmt.Errorf(desc.TextDesc(text.DescKeyErrCryptoReadKey), cause)
}

// InvalidKeySize returns an error when a key file has the wrong size.
//
// Parameters:
//   - got: actual key size in bytes.
//   - want: expected key size in bytes.
//
// Returns:
//   - error: "invalid key size: got N bytes, want M"
func InvalidKeySize(got, want int) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrCryptoInvalidKeySize), got, want,
	)
}

// WriteKey wraps a failure to write a key file.
//
// Parameters:
//   - cause: the underlying write error.
//
// Returns:
//   - error: "write key: <cause>"
func WriteKey(cause error) error {
	return fmt.Errorf(
		desc.TextDesc(text.DescKeyErrCryptoWriteKey), cause,
	)
}
