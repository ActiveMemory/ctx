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
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
	writePad "github.com/ActiveMemory/ctx/internal/write/pad"
)

// ScratchpadPath returns the full path to the scratchpad file.
//
// Returns:
//   - string: Encrypted or plaintext path based on rc.ScratchpadEncrypt()
//   - error: non-nil when the context directory is not declared
func ScratchpadPath() (string, error) {
	ctxDir, err := rc.ContextDir()
	if err != nil {
		return "", err
	}
	if rc.ScratchpadEncrypt() {
		return filepath.Join(ctxDir, pad.Enc), nil
	}
	return filepath.Join(ctxDir, pad.Md), nil
}

// KeyPath returns the full path to the encryption key file.
//
// Triggers legacy key migration on each call, then resolves
// the effective path via rc.KeyPath().
//
// Returns:
//   - string: Resolved key file path
//   - error: propagated from [rc.KeyPath] when the context
//     directory is not declared or otherwise unresolvable
func KeyPath() (string, error) {
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
	kp, kpErr := KeyPath()
	if kpErr != nil {
		return kpErr
	}

	// Key already exists - nothing to do.
	if _, statErr := os.Stat(kp); statErr == nil {
		return nil
	}

	// Encrypted file already exists without a key - we can't generate a new
	// one because it wouldn't decrypt the existing data.
	padPath, padErr := ScratchpadPath()
	if padErr != nil {
		return padErr
	}
	if _, statErr := os.Stat(padPath); statErr == nil {
		return errCrypto.NoKeyAt(kp)
	}

	// First use: generate key.
	key, genErr := crypto.GenerateKey()
	if genErr != nil {
		return errCrypto.GenerateKey(genErr)
	}

	if mkErr := io.SafeMkdirAll(filepath.Dir(kp), fs.PermKeyDir); mkErr != nil {
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
	content, readErr := io.SafeReadUserFile(file.FileGitignore)
	if readErr != nil && !errors.Is(readErr, os.ErrNotExist) {
		return readErr
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
	return io.SafeWriteFile(
		file.FileGitignore,
		[]byte(string(content)+sep+entry+token.NewlineLF), fs.PermFile,
	)
}

// ReadEntriesWithIDs reads the scratchpad and returns
// entries with stable IDs. Auto-migrates legacy entries
// (without ID prefixes) by assigning sequential IDs.
//
// Returns:
//   - []parse.Entry: Entries with stable IDs
//   - error: Non-nil on key or decryption errors
func ReadEntriesWithIDs() ([]parse.Entry, error) {
	data, readErr := readRaw()
	if readErr != nil {
		return nil, readErr
	}
	if data == nil {
		return nil, nil
	}
	return parse.EntriesWithIDs(data), nil
}

// WriteEntriesWithIDs writes ID-prefixed entries to the
// scratchpad file.
//
// Parameters:
//   - cmd: Cobra command for diagnostic output
//   - entries: Entries with stable IDs to write
//
// Returns:
//   - error: Non-nil on key, encryption, or write errors
func WriteEntriesWithIDs(
	cmd *cobra.Command, entries []parse.Entry,
) error {
	path, pathErr := ScratchpadPath()
	if pathErr != nil {
		return pathErr
	}
	plaintext := parse.FormatEntriesWithIDs(entries)

	if !rc.ScratchpadEncrypt() {
		return io.SafeWriteFile(path, plaintext, fs.PermFile)
	}

	if ensureErr := EnsureKey(cmd); ensureErr != nil {
		return ensureErr
	}

	kp, kpErr := KeyPath()
	if kpErr != nil {
		return kpErr
	}
	key, loadErr := crypto.LoadKey(kp)
	if loadErr != nil {
		return errCrypto.LoadKey(loadErr, kp)
	}

	ciphertext, encErr := crypto.Encrypt(key, plaintext)
	if encErr != nil {
		return errCrypto.EncryptFailed(encErr)
	}

	return io.SafeWriteFile(path, ciphertext, fs.PermFile)
}

// ReadEntries reads the scratchpad and returns content
// strings without ID prefixes. IDs are preserved on disk.
// Use ReadEntriesWithIDs when you need stable ID access.
//
// Returns:
//   - []string: Entry content strings (may be empty)
//   - error: Non-nil on key or decryption errors
func ReadEntries() ([]string, error) {
	entries, readErr := ReadEntriesWithIDs()
	if readErr != nil {
		return nil, readErr
	}
	return parse.ToStrings(entries), nil
}

// WriteEntries writes entries to the scratchpad, preserving
// stable IDs. Reads existing IDs from disk, matches entries
// by position, and assigns new IDs for added entries.
//
// Parameters:
//   - cmd: Cobra command for diagnostic output
//   - entries: Content strings to write
//
// Returns:
//   - error: Non-nil on key, encryption, or write errors
func WriteEntries(
	cmd *cobra.Command, entries []string,
) error {
	// Read existing IDs to preserve them.
	existing, _ := ReadEntriesWithIDs()

	idEntries := make([]parse.Entry, len(entries))
	nextID := parse.NextID(existing)

	for i, content := range entries {
		if i < len(existing) {
			idEntries[i] = parse.Entry{
				ID:      existing[i].ID,
				Content: content,
			}
		} else {
			idEntries[i] = parse.Entry{
				ID:      nextID,
				Content: content,
			}
			nextID++
		}
	}

	return WriteEntriesWithIDs(cmd, idEntries)
}
