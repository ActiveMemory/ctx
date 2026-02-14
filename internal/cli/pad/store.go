//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
	"github.com/ActiveMemory/ctx/internal/crypto"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Error messages matching the spec.
const (
	errNoKey       = "encrypted scratchpad found but no key; place your key at .context/.scratchpad.key"
	errDecryptFail = "decryption failed: wrong key?"
	msgEmpty       = "Scratchpad is empty."
	msgKeyCreated  = "Scratchpad key created at %s\nCopy this file to your other machines at the same path.\n"
)

// errEntryRange returns the spec-defined out-of-range error.
func errEntryRange(n, total int) string {
	return fmt.Sprintf("Entry %d does not exist. Scratchpad has %d entries.", n, total)
}

// scratchpadPath returns the full path to the scratchpad file.
func scratchpadPath() string {
	if rc.ScratchpadEncrypt() {
		return filepath.Join(rc.ContextDir(), config.FileScratchpadEnc)
	}
	return filepath.Join(rc.ContextDir(), config.FileScratchpadMd)
}

// keyPath returns the full path to the scratchpad key file.
func keyPath() string {
	return filepath.Join(rc.ContextDir(), config.FileScratchpadKey)
}

// ensureKey generates a scratchpad key if one doesn't exist and there is no
// pre-existing encrypted scratchpad (which would need the original key).
// On first use this lets `ctx pad add` work without requiring `ctx init`.
func ensureKey() error {
	kp := keyPath()

	// Key already exists — nothing to do.
	if _, err := os.Stat(kp); err == nil {
		return nil
	}

	// Encrypted file already exists without a key — we can't generate a new
	// one because it wouldn't decrypt the existing data.
	if _, err := os.Stat(scratchpadPath()); err == nil {
		return errors.New(errNoKey)
	}

	// First use: generate key.
	key, err := crypto.GenerateKey()
	if err != nil {
		return fmt.Errorf("generate scratchpad key: %w", err)
	}

	// Ensure .context/ directory exists.
	if err := os.MkdirAll(filepath.Dir(kp), 0750); err != nil {
		return fmt.Errorf("create context dir: %w", err)
	}

	if err := crypto.SaveKey(kp, key); err != nil {
		return fmt.Errorf("save scratchpad key: %w", err)
	}

	// Best-effort: add key to .gitignore.
	_ = ensureGitignore(rc.ContextDir(), config.FileScratchpadKey)

	fmt.Fprintf(os.Stderr, msgKeyCreated, kp)
	return nil
}

// ensureGitignore adds an entry to .gitignore if not already present.
func ensureGitignore(contextDir, filename string) error {
	entry := filepath.Join(contextDir, filename)
	const gitignorePath = ".gitignore"

	content, err := os.ReadFile(gitignorePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	for _, line := range strings.Split(string(content), "\n") {
		if strings.TrimSpace(line) == entry {
			return nil
		}
	}

	sep := ""
	if len(content) > 0 && !strings.HasSuffix(string(content), "\n") {
		sep = "\n"
	}
	return os.WriteFile(gitignorePath, []byte(string(content)+sep+entry+"\n"), config.PermFile)
}

// readEntries reads the scratchpad and returns its entries.
//
// If the scratchpad file does not exist, it returns an empty slice (no error).
// If the encrypted file exists but the key is missing, it returns an error.
//
// Returns:
//   - []string: The scratchpad entries (may be empty)
//   - error: Non-nil on key or decryption errors
func readEntries() ([]string, error) {
	path := scratchpadPath()

	data, err := os.ReadFile(path) //nolint:gosec // path is constructed from config constants
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read scratchpad: %w", err)
	}

	if !rc.ScratchpadEncrypt() {
		return parseEntries(data), nil
	}

	// Encrypted mode: load key and decrypt
	key, err := crypto.LoadKey(keyPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, errors.New(errNoKey)
		}
		return nil, fmt.Errorf("load key: %w", err)
	}

	plaintext, err := crypto.Decrypt(key, data)
	if err != nil {
		return nil, errors.New(errDecryptFail)
	}

	return parseEntries(plaintext), nil
}

// writeEntries writes entries to the scratchpad file.
//
// In encrypted mode, the entries are encrypted with AES-256-GCM before
// writing. In plaintext mode, they are written as a newline-delimited file.
//
// Parameters:
//   - entries: The scratchpad entries to write
//
// Returns:
//   - error: Non-nil on key, encryption, or file write errors
func writeEntries(entries []string) error {
	path := scratchpadPath()
	plaintext := formatEntries(entries)

	if !rc.ScratchpadEncrypt() {
		return os.WriteFile(path, plaintext, config.PermFile)
	}

	// Encrypted mode: ensure key exists (auto-generate on first use).
	if err := ensureKey(); err != nil {
		return err
	}

	key, err := crypto.LoadKey(keyPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return errors.New(errNoKey)
		}
		return fmt.Errorf("load key: %w", err)
	}

	ciphertext, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		return fmt.Errorf("encrypt: %w", err)
	}

	return os.WriteFile(path, ciphertext, config.PermFile)
}

// parseEntries splits raw bytes into entry lines, filtering empty lines.
func parseEntries(data []byte) []string {
	if len(data) == 0 {
		return nil
	}
	lines := strings.Split(string(data), "\n")
	var entries []string
	for _, line := range lines {
		if line != "" {
			entries = append(entries, line)
		}
	}
	return entries
}

// formatEntries joins entries with newlines and adds a trailing newline.
func formatEntries(entries []string) []byte {
	if len(entries) == 0 {
		return nil
	}
	return []byte(strings.Join(entries, "\n") + "\n")
}
