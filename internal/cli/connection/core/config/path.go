//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

import (
	"path/filepath"

	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
	"github.com/ActiveMemory/ctx/internal/crypto"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// filePath returns the full path to .connect.enc.
//
// Returns:
//   - string: Absolute path to the encrypted connect file
func filePath() string {
	return filepath.Join(rc.ContextDir(), cfgHub.FileConnect)
}

// loadKey reads the encryption key from the global key
// path.
//
// Returns:
//   - []byte: raw encryption key bytes
//   - error: non-nil if the key file cannot be read
func loadKey() ([]byte, error) {
	return crypto.LoadKey(crypto.GlobalKeyPath())
}
