//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/crypto"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// filePath returns the full path to .connect.enc.
func filePath() string {
	return filepath.Join(rc.ContextDir(), connectFile)
}

// loadKey reads the encryption key.
func loadKey() ([]byte, error) {
	return crypto.LoadKey(crypto.GlobalKeyPath())
}
