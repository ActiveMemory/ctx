//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package state manages journal processing state via an external JSON file.
//
// Instead of embedding markers (<!-- normalized: ... -->) inside journal
// files, which causes false positives when journal content includes those
// exact strings, state is tracked in .context/journal/.state.json.
package state

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/file"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/journal"
)

// CurrentVersion is the schema version for the state file.
const CurrentVersion = 1

// Load reads the state file from the journal directory.
//
// If the file does not exist, an empty state is returned (not an error).
//
// Parameters:
//   - journalDir: path to the journal directory
//
// Returns:
//   - *JournalState: loaded or empty state
//   - error: non-nil if the file exists but cannot be read or parsed
func Load(journalDir string) (*JournalState, error) {
	path := filepath.Join(journalDir, journal.FileState)

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
//
// Parameters:
//   - journalDir: path to the journal directory
//
// Returns:
//   - error: non-nil if marshalling or file write fails
func (s *JournalState) Save(journalDir string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	path := filepath.Join(journalDir, journal.FileState)
	tmp := path + ".tmp"

	if err := os.WriteFile(tmp, data, fs.PermFile); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// MarkImported records that a file was imported.
//
// Parameters:
//   - filename: journal entry filename (e.g., "2026-01-21-session.md")
func (s *JournalState) MarkImported(filename string) {
	ff := s.Entries[filename]
	ff.Exported = today()
	s.Entries[filename] = ff
}

// MarkEnriched records that a file was enriched.
//
// Parameters:
//   - filename: journal entry filename
func (s *JournalState) MarkEnriched(filename string) {
	ff := s.Entries[filename]
	ff.Enriched = today()
	s.Entries[filename] = ff
}

// MarkNormalized records that a file was normalized.
//
// Parameters:
//   - filename: journal entry filename
func (s *JournalState) MarkNormalized(filename string) {
	ff := s.Entries[filename]
	ff.Normalized = today()
	s.Entries[filename] = ff
}

// MarkFencesVerified records that a file's fences were verified.
//
// Parameters:
//   - filename: journal entry filename
func (s *JournalState) MarkFencesVerified(filename string) {
	ff := s.Entries[filename]
	ff.FencesVerified = today()
	s.Entries[filename] = ff
}

// Mark sets an arbitrary stage to today's date.
//
// Parameters:
//   - filename: journal entry filename (e.g., "2026-01-21-session.md")
//   - stage: one of ValidStages (exported, enriched, normalized,
//     fences_verified, locked)
//
// Returns:
//   - bool: false if stage is not recognized
func (s *JournalState) Mark(filename, stage string) bool {
	ff := s.Entries[filename]
	switch stage {
	case journal.StageExported:
		ff.Exported = today()
	case journal.StageEnriched:
		ff.Enriched = today()
	case journal.StageNormalized:
		ff.Normalized = today()
	case journal.StageFencesVerified:
		ff.FencesVerified = today()
	case journal.StageLocked:
		ff.Locked = today()
	default:
		return false
	}
	s.Entries[filename] = ff
	return true
}

// Clear removes a stage value, resetting it to empty.
//
// Parameters:
//   - filename: journal entry filename
//   - stage: one of ValidStages
//
// Returns:
//   - bool: false if stage is not recognized
func (s *JournalState) Clear(filename, stage string) bool {
	ff := s.Entries[filename]
	switch stage {
	case journal.StageExported:
		ff.Exported = ""
	case journal.StageEnriched:
		ff.Enriched = ""
	case journal.StageNormalized:
		ff.Normalized = ""
	case journal.StageFencesVerified:
		ff.FencesVerified = ""
	case journal.StageLocked:
		ff.Locked = ""
	default:
		return false
	}
	s.Entries[filename] = ff
	return true
}

// Locked reports whether the file is protected from export regeneration.
//
// Parameters:
//   - filename: journal entry filename
//
// Returns:
//   - bool: true if the file has a lock date recorded
func (s *JournalState) Locked(filename string) bool {
	return s.Entries[filename].Locked != ""
}

// Rename moves state from an old filename to a new one, preserving all
// fields. If old does not exist in state, this is a no-op.
//
// Parameters:
//   - oldName: current filename in state
//   - newName: target filename
func (s *JournalState) Rename(oldName, newName string) {
	ff, ok := s.Entries[oldName]
	if !ok {
		return
	}
	s.Entries[newName] = ff
	delete(s.Entries, oldName)
}

// ClearEnriched removes the enriched date for a file, resetting it to
// unenriched. Used when --force re-export discards frontmatter.
//
// Parameters:
//   - filename: journal entry filename
func (s *JournalState) ClearEnriched(filename string) {
	ff := s.Entries[filename]
	ff.Enriched = ""
	s.Entries[filename] = ff
}

// Enriched reports whether the file has been enriched.
//
// Parameters:
//   - filename: journal entry filename
//
// Returns:
//   - bool: true if the file has an enriched date
func (s *JournalState) Enriched(filename string) bool {
	return s.Entries[filename].Enriched != ""
}

// Normalized reports whether the file has been normalized.
//
// Parameters:
//   - filename: journal entry filename
//
// Returns:
//   - bool: true if the file has a normalized date
func (s *JournalState) Normalized(filename string) bool {
	return s.Entries[filename].Normalized != ""
}

// FencesVerified reports whether the file's fences have been verified.
//
// Parameters:
//   - filename: journal entry filename
//
// Returns:
//   - bool: true if the file has a fences-verified date
func (s *JournalState) FencesVerified(filename string) bool {
	return s.Entries[filename].FencesVerified != ""
}

// Exported reports whether the file has been exported.
//
// Parameters:
//   - filename: journal entry filename
//
// Returns:
//   - bool: true if the file has an exported date
func (s *JournalState) Exported(filename string) bool {
	return s.Entries[filename].Exported != ""
}

// CountUnenriched counts .md files in the directory that lack an
// enriched date in the state file.
//
// Parameters:
//   - journalDir: path to the journal directory
//
// Returns:
//   - int: number of unenriched Markdown files
func (s *JournalState) CountUnenriched(journalDir string) int {
	entries, err := os.ReadDir(journalDir)
	if err != nil {
		return 0
	}

	count := 0
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != file.ExtMarkdown {
			continue
		}
		if !s.Enriched(entry.Name()) {
			count++
		}
	}
	return count
}

// ValidStages lists the recognized stage names for Mark and Clear.
var ValidStages = []string{
	journal.StageExported,
	journal.StageEnriched,
	journal.StageNormalized,
	journal.StageFencesVerified,
	journal.StageLocked,
}
