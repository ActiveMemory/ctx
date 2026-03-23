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
	"strings"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/pad/core/parse"
	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/pad"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/crypto"
	errCrypto "github.com/ActiveMemory/ctx/internal/err/crypto"
	errPad "github.com/ActiveMemory/ctx/internal/err/pad"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
	writePad "github.com/ActiveMemory/ctx/internal/write/pad"
)

// ScratchpadPath returns the full path to the scratchpad file.
//
// Returns:
//   - string: Encrypted or plaintext path based on rc.ScratchpadEncrypt()
func ScratchpadPath() string {
	if rc.ScratchpadEncrypt() {
		return filepath.Join(rc.ContextDir(), pad.Enc)
	}
	return filepath.Join(rc.ContextDir(), pad.Md)
}

// KeyPath returns the full path to the encryption key file.
//
// Triggers legacy key migration on each call, then resolves
// the effective path via rc.KeyPath().
//
// Returns:
//   - string: Resolved key file path
func KeyPath() string {
	return rc.KeyPath()
}

// EnsureKey generates a scratchpad key when none exists.
//
// If an encrypted scratchpad already exists without a key, returns an
// error (a new key would not decrypt the existing data). On first use
// this lets `ctx pad add` work without requiring `ctx init`.
//
// Parameters:
//   - cmd: Cobra command for diagnostic output
//
// Returns:
//   - error: Non-nil on missing key with existing data, or generation failure
func EnsureKey(cmd *cobra.Command) error {
	kp := KeyPath()

	// Key already exists - nothing to do.
	if _, err := os.Stat(kp); err == nil {
		return nil
	}

	// Encrypted file already exists without a key - we can't generate a new
	// one because it wouldn't decrypt the existing data.
	if _, err := os.Stat(ScratchpadPath()); err == nil {
		return errCrypto.NoKeyAt(kp)
	}

	// First use: generate key.
	key, genErr := crypto.GenerateKey()
	if genErr != nil {
		return errCrypto.GenerateKey(genErr)
	}

	if mkErr := os.MkdirAll(filepath.Dir(kp), fs.PermKeyDir); mkErr != nil {
		return errCrypto.MkdirKeyDir(mkErr)
	}

	if saveErr := crypto.SaveKey(kp, key); saveErr != nil {
		return errCrypto.SaveKey(saveErr)
	}

	writePad.KeyCreated(cmd, kp)
	return nil
}

// EnsureGitignore adds an entry to .gitignore if not already present.
//
// Parameters:
//   - contextDir: The .context directory path
//   - filename: The file to add (joined with contextDir)
//
// Returns:
//   - error: Non-nil on read/write failure
func EnsureGitignore(contextDir, filename string) error {
	entry := filepath.Join(contextDir, filename)
	content, err := os.ReadFile(file.FileGitignore)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	for _, line := range strings.Split(string(content), token.NewlineLF) {
		if strings.TrimSpace(line) == entry {
			return nil
		}
	}

	sep := ""
	if len(content) > 0 && !strings.HasSuffix(string(content), token.NewlineLF) {
		sep = token.NewlineLF
	}
	return os.WriteFile(
		file.FileGitignore,
		[]byte(string(content)+sep+entry+token.NewlineLF), fs.PermFile,
	)
}

// ReadEntries reads the scratchpad and returns its entries.
//
// If the scratchpad file does not exist, it returns an empty slice (no error).
// If the encrypted file exists but the key is missing, it returns an error.
//
// Returns:
//   - []string: The scratchpad entries (may be empty)
//   - error: Non-nil on key or decryption errors
func ReadEntries() ([]string, error) {
	path := ScratchpadPath()
	dir := filepath.Dir(path)
	name := filepath.Base(path)

	data, err := io.SafeReadFile(dir, name)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, errPad.ReadScratchpad(err)
	}

	if !rc.ScratchpadEncrypt() {
		return parse.ParseEntries(data), nil
	}

	kp := KeyPath()
	key, loadErr := crypto.LoadKey(kp)
	if loadErr != nil {
		return nil, errCrypto.LoadKey(loadErr, kp)
	}

	plaintext, decErr := crypto.Decrypt(key, data)
	if decErr != nil {
		return nil, errCrypto.DecryptFailed()
	}

	return parse.ParseEntries(plaintext), nil
}

// WriteEntries writes entries to the scratchpad file.
//
// In encrypted mode, the entries are encrypted with AES-256-GCM before
// writing. In plaintext mode, they are written as a newline-delimited file.
//
// Parameters:
//   - cmd: Cobra command for diagnostic output
//   - entries: The scratchpad entries to write
//
// Returns:
//   - error: Non-nil on key, encryption, or file write errors
func WriteEntries(cmd *cobra.Command, entries []string) error {
	path := ScratchpadPath()
	plaintext := parse.FormatEntries(entries)

	if !rc.ScratchpadEncrypt() {
		return os.WriteFile(path, plaintext, fs.PermFile)
	}

	if err := EnsureKey(cmd); err != nil {
		return err
	}

	kp := KeyPath()
	key, loadErr := crypto.LoadKey(kp)
	if loadErr != nil {
		return errCrypto.LoadKey(loadErr, kp)
	}

	ciphertext, encErr := crypto.Encrypt(key, plaintext)
	if encErr != nil {
		return errCrypto.EncryptFailed(encErr)
	}

	return os.WriteFile(path, ciphertext, fs.PermFile)
}
