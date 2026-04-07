//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package register

import (
	"encoding/json"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/crypto"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// connectFile is the encrypted connection config filename.
const connectFile = ".connect.enc"

// saveConfig encrypts and writes the connection config.
func saveConfig(hubAddr, bearerToken string) error {
	cfg := connConfig{
		HubAddr: hubAddr,
		Token:   bearerToken,
	}
	data, marshalErr := json.Marshal(cfg)
	if marshalErr != nil {
		return marshalErr
	}

	key, keyErr := crypto.LoadKey(
		crypto.GlobalKeyPath(),
	)
	if keyErr != nil {
		return keyErr
	}

	encrypted, encErr := crypto.Encrypt(key, data)
	if encErr != nil {
		return encErr
	}

	configPath := filepath.Join(
		rc.ContextDir(), connectFile,
	)
	return io.SafeWriteFile(
		configPath, encrypted, fs.PermSecret,
	)
}
