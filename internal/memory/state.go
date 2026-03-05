//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
)

// State tracks memory bridge sync timestamps and (in future phases)
// import/publish progress.
type State struct {
	LastSync       *time.Time `json:"last_sync"`
	LastImport     *time.Time `json:"last_import"`
	LastPublish    *time.Time `json:"last_publish"`
	ImportedHashes []string   `json:"imported_hashes"`
}

// LoadState reads the sync state from .context/state/memory-import.json.
// Returns a zero-value State if the file does not exist.
func LoadState(contextDir string) (State, error) {
	path := statePath(contextDir)
	data, readErr := os.ReadFile(path) //nolint:gosec // project-local state path
	if readErr != nil {
		if errors.Is(readErr, os.ErrNotExist) {
			return State{ImportedHashes: []string{}}, nil
		}
		return State{}, readErr
	}

	var s State
	if unmarshalErr := json.Unmarshal(data, &s); unmarshalErr != nil {
		return State{}, unmarshalErr
	}
	if s.ImportedHashes == nil {
		s.ImportedHashes = []string{}
	}
	return s, nil
}

// SaveState writes the sync state to .context/state/memory-import.json.
func SaveState(contextDir string, s State) error {
	path := statePath(contextDir)
	dir := filepath.Dir(path)
	if mkErr := os.MkdirAll(dir, config.PermExec); mkErr != nil {
		return mkErr
	}

	data, marshalErr := json.MarshalIndent(s, "", "  ")
	if marshalErr != nil {
		return marshalErr
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, config.PermFile)
}

// MarkSynced updates the state with the current timestamp.
func (s *State) MarkSynced() {
	now := time.Now().UTC()
	s.LastSync = &now
}

// EntryHash computes a deduplication hash for an entry.
// Uses SHA-256 of the text, truncated to 16 hex chars.
func EntryHash(text string) string {
	h := sha256.Sum256([]byte(text))
	return fmt.Sprintf("%x", h[:8])
}

// Imported reports whether an entry hash has already been imported.
// Stored entries use format "hash:target:date"; matches on hash prefix.
func (s *State) Imported(hash string) bool {
	prefix := hash + ":"
	for _, h := range s.ImportedHashes {
		if h == hash || len(h) > len(hash) && h[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

// MarkImported records an entry hash with its target and date.
func (s *State) MarkImported(hash, target string) {
	date := time.Now().Format("2006-01-02")
	entry := fmt.Sprintf("%s:%s:%s", hash, target, date)
	s.ImportedHashes = append(s.ImportedHashes, entry)
}

// MarkImportedDone updates LastImport to the current time.
func (s *State) MarkImportedDone() {
	now := time.Now().UTC()
	s.LastImport = &now
}

func statePath(contextDir string) string {
	return filepath.Join(contextDir, config.DirState, config.FileMemoryState)
}
