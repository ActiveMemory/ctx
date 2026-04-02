//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	cfgPad "github.com/ActiveMemory/ctx/internal/config/pad"
	"github.com/ActiveMemory/ctx/internal/crypto"
	errCrypto "github.com/ActiveMemory/ctx/internal/err/crypto"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write/initialize"
)

// Setup configures the scratchpad key or plaintext file.
//
// When encryption is enabled (default):
//   - Generates a 256-bit key at ~/.ctx/ if not present
//   - Warns if .enc exists but no key
//
// When encryption is disabled:
//   - Creates empty .context/scratchpad.md if not present
//
// Parameters:
//   - cmd: Cobra command for output
//   - contextDir: The .context/ directory path
//
// Returns:
//   - error: Non-nil if key generation or file operations fail
func Setup(cmd *cobra.Command, contextDir string) error {
	if !rc.ScratchpadEncrypt() {
		return setupPlaintext(cmd, contextDir)
	}
	return setupEncrypted(cmd, contextDir)
}

func setupPlaintext(cmd *cobra.Command, contextDir string) error {
	mdPath := filepath.Join(contextDir, cfgPad.Md)
	if _, statErr := os.Stat(mdPath); statErr != nil {
		if writeErr := ctxIo.SafeWriteFile(mdPath, nil, cfgFs.PermFile); writeErr != nil {
			return writeErr
		}
		initialize.InfoScratchpadPlaintext(cmd, mdPath)
	} else {
		initialize.InfoExistsSkipped(cmd, mdPath)
	}
	return nil
}

func setupEncrypted(cmd *cobra.Command, contextDir string) error {
	kPath := rc.KeyPath()
	encPath := filepath.Join(contextDir, cfgPad.Enc)

	// Check if the key already exists (idempotent)
	if _, keyStatErr := os.Stat(kPath); keyStatErr == nil {
		initialize.InfoExistsSkipped(cmd, kPath)
		return nil
	}

	// Warn if the encrypted file exists but no key
	if _, encStatErr := os.Stat(encPath); encStatErr == nil {
		initialize.InfoScratchpadNoKey(cmd, kPath)
		return nil
	}

	// Ensure the key directory exists.
	if mkdirErr := ctxIo.SafeMkdirAll(
		filepath.Dir(kPath), cfgFs.PermKeyDir,
	); mkdirErr != nil {
		return errCrypto.MkdirKeyDir(mkdirErr)
	}

	// Generate key
	key, genErr := crypto.GenerateKey()
	if genErr != nil {
		return errCrypto.GenerateKey(genErr)
	}

	if saveErr := crypto.SaveKey(kPath, key); saveErr != nil {
		return errCrypto.SaveKey(saveErr)
	}
	initialize.InfoScratchpadKeyCreated(cmd, kPath)

	return nil
}
