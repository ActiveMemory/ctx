//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package err

import (
	"errors"
	"fmt"
	"os"

	"github.com/ActiveMemory/ctx/internal/assets"
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
		assets.TextDesc(assets.TextDescKeyErrCryptoLoadKey), cause,
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
		assets.TextDesc(assets.TextDescKeyErrCryptoEncryptFailed), cause,
	)
}

// DecryptFailed returns an error indicating decryption failure.
//
// Returns:
//   - error: "decryption failed: wrong key?"
func DecryptFailed() error {
	return errors.New(assets.TextDesc(assets.TextDescKeyErrCryptoDecryptFailed))
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
		assets.TextDesc(assets.TextDescKeyErrCryptoNoKeyAt), path,
	)
}

// GenerateKey wraps a failure to generate an encryption key.
//
// Parameters:
//   - cause: the underlying error from key generation
//
// Returns:
//   - error: "failed to generate scratchpad key: <cause>"
func GenerateKey(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrCryptoGenerateScratchpadKey), cause,
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
		assets.TextDesc(assets.TextDescKeyErrCryptoSaveKey), cause)
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
		assets.TextDesc(assets.TextDescKeyErrCryptoMkdirKeyDir), cause,
	)
}

// CryptoCreateCipher wraps a failure to create an AES cipher.
//
// Parameters:
//   - cause: the underlying crypto error.
//
// Returns:
//   - error: "create cipher: <cause>"
func CryptoCreateCipher(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrCryptoCreateCipher), cause,
	)
}

// CryptoCreateGCM wraps a failure to create a GCM instance.
//
// Parameters:
//   - cause: the underlying crypto error.
//
// Returns:
//   - error: "create GCM: <cause>"
func CryptoCreateGCM(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrCryptoCreateGCM), cause,
	)
}

// CryptoGenerateNonce wraps a failure to generate a random nonce.
//
// Parameters:
//   - cause: the underlying IO error.
//
// Returns:
//   - error: "generate nonce: <cause>"
func CryptoGenerateNonce(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrCryptoGenerateNonce), cause,
	)
}

// CryptoGenerateKey wraps a failure to generate a random key.
//
// Parameters:
//   - cause: the underlying IO error.
//
// Returns:
//   - error: "generate key: <cause>"
func CryptoGenerateKey(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrCryptoGenerateKey), cause,
	)
}

// CryptoCiphertextTooShort returns an error when ciphertext is shorter
// than the nonce size.
//
// Returns:
//   - error: "ciphertext too short"
func CryptoCiphertextTooShort() error {
	return errors.New(
		assets.TextDesc(assets.TextDescKeyErrCryptoCiphertextTooShort),
	)
}

// CryptoDecrypt wraps a decryption failure with cause.
//
// Parameters:
//   - cause: the underlying decryption error.
//
// Returns:
//   - error: "decrypt: <cause>"
func CryptoDecrypt(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrCryptoDecrypt), cause,
	)
}

// CryptoReadKey wraps a failure to read a key file.
//
// Parameters:
//   - cause: the underlying read error.
//
// Returns:
//   - error: "read key: <cause>"
func CryptoReadKey(cause error) error {
	return fmt.Errorf(assets.TextDesc(assets.TextDescKeyErrCryptoReadKey), cause)
}

// CryptoInvalidKeySize returns an error when a key file has the wrong size.
//
// Parameters:
//   - got: actual key size in bytes.
//   - want: expected key size in bytes.
//
// Returns:
//   - error: "invalid key size: got N bytes, want M"
func CryptoInvalidKeySize(got, want int) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrCryptoInvalidKeySize), got, want,
	)
}

// CryptoWriteKey wraps a failure to write a key file.
//
// Parameters:
//   - cause: the underlying write error.
//
// Returns:
//   - error: "write key: <cause>"
func CryptoWriteKey(cause error) error {
	return fmt.Errorf(
		assets.TextDesc(assets.TextDescKeyErrCryptoWriteKey), cause,
	)
}
