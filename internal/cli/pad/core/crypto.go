//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"github.com/ActiveMemory/ctx/internal/crypto"
	ctxErr "github.com/ActiveMemory/ctx/internal/err/crypto"
	"github.com/ActiveMemory/ctx/internal/io"
)

// DecryptFile reads and decrypts a single file, returning its entries.
//
// Parameters:
//   - key: AES-256 encryption key
//   - baseDir: Directory containing the encrypted file
//   - filename: Name of the encrypted file
//
// Returns:
//   - []string: Decrypted entries
//   - error: Non-nil on read or decryption failure
func DecryptFile(key []byte, baseDir, filename string) ([]string, error) {
	data, readErr := io.SafeReadFile(baseDir, filename)
	if readErr != nil {
		return nil, readErr
	}

	plaintext, decErr := crypto.Decrypt(key, data)
	if decErr != nil {
		return nil, ctxErr.DecryptFailed()
	}

	return ParseEntries(plaintext), nil
}
