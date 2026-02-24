//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package state manages journal processing state via an external JSON file.
//
// Instead of embedding markers (<!-- normalized: ... -->) inside journal
// files — which causes false positives when journal content includes those
// exact strings — state is tracked in .context/journal/.state.json.
package state

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/config"
)

// CurrentVersion is the schema version for the state file.
const CurrentVersion = 1

// JournalState is the top-level state file structure.
type JournalState struct {
	Version int                  `json:"version"`
	Entries map[string]FileState `json:"entries"`
}

// FileState tracks processing stages for a single journal entry.
// Values are date strings (YYYY-MM-DD) indicating when the stage completed.
type FileState struct {
	Exported       string `json:"exported,omitempty"`
	Enriched       string `json:"enriched,omitempty"`
	Normalized     string `json:"normalized,omitempty"`
	FencesVerified string `json:"fences_verified,omitempty"`
	Locked         string `json:"locked,omitempty"`
}

// Load reads the state file from the journal directory. If the file does
// not exist, an empty state is returned (not an error).
func Load(journalDir string) (*JournalState, error) {
	path := filepath.Join(journalDir, config.FileJournalState)

	data, err := os.ReadFile(filepath.Clean(path))
	if os.IsNotExist(err) {
		return &JournalState{
			Version: CurrentVersion,
			Entries: make(map[string]FileState),
		}, nil
	}
	if err != nil {
		return nil, err
	}

	var s JournalState
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	if s.Entries == nil {
		s.Entries = make(map[string]FileState)
	}
	return &s, nil
}

// Save writes the state file atomically (temp + rename) to the journal
// directory.
func (s *JournalState) Save(journalDir string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	path := filepath.Join(journalDir, config.FileJournalState)
	tmp := path + ".tmp"

	if err := os.WriteFile(tmp, data, config.PermFile); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// today returns today's date as YYYY-MM-DD.
func today() string {
	return time.Now().Format("2006-01-02")
}

// MarkExported records that a file was exported.
func (s *JournalState) MarkExported(filename string) {
	fs := s.Entries[filename]
	fs.Exported = today()
	s.Entries[filename] = fs
}

// MarkEnriched records that a file was enriched.
func (s *JournalState) MarkEnriched(filename string) {
	fs := s.Entries[filename]
	fs.Enriched = today()
	s.Entries[filename] = fs
}

// MarkNormalized records that a file was normalized.
func (s *JournalState) MarkNormalized(filename string) {
	fs := s.Entries[filename]
	fs.Normalized = today()
	s.Entries[filename] = fs
}

// MarkFencesVerified records that a file's fences were verified.
func (s *JournalState) MarkFencesVerified(filename string) {
	fs := s.Entries[filename]
	fs.FencesVerified = today()
	s.Entries[filename] = fs
}

// Mark sets an arbitrary stage to today's date.
//
// Parameters:
//   - filename: Journal entry filename (e.g. "2026-01-21-session.md")
//   - stage: One of ValidStages (exported, enriched, normalized,
//     fences_verified, locked)
//
// Returns:
//   - bool: False if stage is not recognized
func (s *JournalState) Mark(filename, stage string) bool {
	fs := s.Entries[filename]
	switch stage {
	case "exported":
		fs.Exported = today()
	case "enriched":
		fs.Enriched = today()
	case "normalized":
		fs.Normalized = today()
	case "fences_verified":
		fs.FencesVerified = today()
	case "locked":
		fs.Locked = today()
	default:
		return false
	}
	s.Entries[filename] = fs
	return true
}

// Clear removes a stage value, resetting it to empty.
//
// Parameters:
//   - filename: Journal entry filename
//   - stage: One of ValidStages
//
// Returns:
//   - bool: False if stage is not recognized
func (s *JournalState) Clear(filename, stage string) bool {
	fs := s.Entries[filename]
	switch stage {
	case "exported":
		fs.Exported = ""
	case "enriched":
		fs.Enriched = ""
	case "normalized":
		fs.Normalized = ""
	case "fences_verified":
		fs.FencesVerified = ""
	case "locked":
		fs.Locked = ""
	default:
		return false
	}
	s.Entries[filename] = fs
	return true
}

// Locked reports whether the file is protected from export regeneration.
//
// Parameters:
//   - filename: Journal entry filename
//
// Returns:
//   - bool: True if the file has a lock date recorded
func (s *JournalState) Locked(filename string) bool {
	return s.Entries[filename].Locked != ""
}

// Rename moves state from an old filename to a new one, preserving all
// fields. If old does not exist in state, this is a no-op.
func (s *JournalState) Rename(oldName, newName string) {
	fs, ok := s.Entries[oldName]
	if !ok {
		return
	}
	s.Entries[newName] = fs
	delete(s.Entries, oldName)
}

// ClearEnriched removes the enriched date for a file, resetting it to
// unenriched. This is used when --force re-export discards frontmatter.
func (s *JournalState) ClearEnriched(filename string) {
	fs := s.Entries[filename]
	fs.Enriched = ""
	s.Entries[filename] = fs
}

// IsEnriched reports whether the file has been enriched.
func (s *JournalState) IsEnriched(filename string) bool {
	return s.Entries[filename].Enriched != ""
}

// IsNormalized reports whether the file has been normalized.
func (s *JournalState) IsNormalized(filename string) bool {
	return s.Entries[filename].Normalized != ""
}

// IsFencesVerified reports whether the file's fences have been verified.
func (s *JournalState) IsFencesVerified(filename string) bool {
	return s.Entries[filename].FencesVerified != ""
}

// IsExported reports whether the file has been exported.
func (s *JournalState) IsExported(filename string) bool {
	return s.Entries[filename].Exported != ""
}

// CountUnenriched counts .md files in the directory that lack an enriched
// date in the state file.
func (s *JournalState) CountUnenriched(journalDir string) int {
	entries, err := os.ReadDir(journalDir)
	if err != nil {
		return 0
	}

	count := 0
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != config.ExtMarkdown {
			continue
		}
		if !s.IsEnriched(entry.Name()) {
			count++
		}
	}
	return count
}

// ValidStages lists the recognized stage names for Mark() and Clear().
var ValidStages = []string{
	"exported", "enriched", "normalized", "fences_verified", "locked",
}
