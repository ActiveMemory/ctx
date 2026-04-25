//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package config

import (
	"encoding/json"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/crypto"
	"github.com/ActiveMemory/ctx/internal/io"
)

// Save encrypts and writes the connection config.
//
// Parameters:
//   - cfg: connection config to persist
//
// Returns:
//   - error: non-nil if encryption or write fails
func Save(cfg Config) error {
	data, marshalErr := json.Marshal(cfg)
	if marshalErr != nil {
		return marshalErr
	}

	key, keyErr := loadKey()
	if keyErr != nil {
		return keyErr
	}

	encrypted, encErr := crypto.Encrypt(key, data)
	if encErr != nil {
		return encErr
	}

	path, pathErr := filePath()
	if pathErr != nil {
		return pathErr
	}
	return io.SafeWriteFile(
		path, encrypted, fs.PermSecret,
	)
}

// Load reads and decrypts the connection config.
//
// Returns:
//   - Config: decrypted connection config
//   - error: non-nil if file missing or decryption fails
func Load() (Config, error) {
	var cfg Config

	path, pathErr := filePath()
	if pathErr != nil {
		return cfg, pathErr
	}
	encrypted, readErr := io.SafeReadUserFile(path)
	if readErr != nil {
		return cfg, readErr
	}

	key, keyErr := loadKey()
	if keyErr != nil {
		return cfg, keyErr
	}

	data, decErr := crypto.Decrypt(key, encrypted)
	if decErr != nil {
		return cfg, decErr
	}

	return cfg, json.Unmarshal(data, &cfg)
}
