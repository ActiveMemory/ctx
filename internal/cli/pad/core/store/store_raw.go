//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package store

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/crypto"
	errCrypto "github.com/ActiveMemory/ctx/internal/err/crypto"
	errPad "github.com/ActiveMemory/ctx/internal/err/pad"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// readRaw reads and decrypts the raw scratchpad content.
//
// Returns:
//   - []byte: Decrypted plaintext, or nil if file missing
//   - error: Non-nil on key or decryption errors
func readRaw() ([]byte, error) {
	path, pathErr := ScratchpadPath()
	if pathErr != nil {
		return nil, pathErr
	}
	dir := filepath.Dir(path)
	name := filepath.Base(path)

	data, readErr := io.SafeReadFile(dir, name)
	if readErr != nil {
		if errors.Is(readErr, os.ErrNotExist) {
			return nil, nil
		}
		return nil, errPad.Read(readErr)
	}

	if !rc.ScratchpadEncrypt() {
		return data, nil
	}

	kp, kpErr := KeyPath()
	if kpErr != nil {
		return nil, kpErr
	}
	key, loadErr := crypto.LoadKey(kp)
	if loadErr != nil {
		return nil, errCrypto.LoadKey(loadErr, kp)
	}

	plaintext, decErr := crypto.Decrypt(key, data)
	if decErr != nil {
		return nil, errCrypto.DecryptFailed()
	}

	return plaintext, nil
}
