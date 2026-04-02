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
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/format"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
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
//   - *State: loaded or empty state
//   - error: non-nil if the file exists but cannot be read or parsed
func Load(journalDir string) (*State, error) {
	path := filepath.Join(journalDir, journal.File)

	data, readErr := ctxIo.SafeReadUserFile(filepath.Clean(path))
	if os.IsNotExist(readErr) {
		return &State{
			Version: CurrentVersion,
			Entries: make(map[string]File),
		}, nil
	}
	if readErr != nil {
		return nil, readErr
	}

	var s State
	if unmarshalErr := json.Unmarshal(data, &s); unmarshalErr != nil {
		return nil, unmarshalErr
	}
	if s.Entries == nil {
		s.Entries = make(map[string]File)
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
func (s *State) Save(journalDir string) error {
	data, marshalErr := json.MarshalIndent(s, "", "  ")
	if marshalErr != nil {
		return marshalErr
	}
	data = append(data, token.NewlineLF[0])

	path := filepath.Join(journalDir, journal.File)
	tmp := path + file.ExtTmp

	if writeErr := ctxIo.SafeWriteFile(tmp, data, fs.PermFile); writeErr != nil {
		return writeErr
	}
	return os.Rename(tmp, path)
}

// MarkImported records that a file was imported.
//
// Parameters:
//   - filename: journal entry filename (e.g., "2026-01-21-session.md")
func (s *State) MarkImported(filename string) {
	ff := s.Entries[filename]
	ff.Exported = format.Today()
	s.Entries[filename] = ff
}

// MarkEnriched records that a file was enriched.
//
// Parameters:
//   - filename: journal entry filename
func (s *State) MarkEnriched(filename string) {
	ff := s.Entries[filename]
	ff.Enriched = format.Today()
	s.Entries[filename] = ff
}

// MarkNormalized records that a file was normalized.
//
// Parameters:
//   - filename: journal entry filename
func (s *State) MarkNormalized(filename string) {
	ff := s.Entries[filename]
	ff.Normalized = format.Today()
	s.Entries[filename] = ff
}

// MarkFencesVerified records that a file's fences were verified.
//
// Parameters:
//   - filename: journal entry filename
func (s *State) MarkFencesVerified(filename string) {
	ff := s.Entries[filename]
	ff.FencesVerified = format.Today()
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
func (s *State) Mark(filename, stage string) bool {
	ff := s.Entries[filename]
	switch stage {
	case journal.StageExported:
		ff.Exported = format.Today()
	case journal.StageEnriched:
		ff.Enriched = format.Today()
	case journal.StageNormalized:
		ff.Normalized = format.Today()
	case journal.StageFencesVerified:
		ff.FencesVerified = format.Today()
	case journal.StageLocked:
		ff.Locked = format.Today()
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
func (s *State) Clear(filename, stage string) bool {
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
func (s *State) Locked(filename string) bool {
	return s.Entries[filename].Locked != ""
}

// Rename moves state from an old filename to a new one, preserving all
// fields. If old does not exist in state, this is a no-op.
//
// Parameters:
//   - oldName: current filename in state
//   - newName: target filename
func (s *State) Rename(oldName, newName string) {
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
func (s *State) ClearEnriched(filename string) {
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
func (s *State) Enriched(filename string) bool {
	return s.Entries[filename].Enriched != ""
}

// Normalized reports whether the file has been normalized.
//
// Parameters:
//   - filename: journal entry filename
//
// Returns:
//   - bool: true if the file has a normalized date
func (s *State) Normalized(filename string) bool {
	return s.Entries[filename].Normalized != ""
}

// FencesVerified reports whether the file's fences have been verified.
//
// Parameters:
//   - filename: journal entry filename
//
// Returns:
//   - bool: true if the file has a fences-verified date
func (s *State) FencesVerified(filename string) bool {
	return s.Entries[filename].FencesVerified != ""
}

// Exported reports whether the file has been exported.
//
// Parameters:
//   - filename: journal entry filename
//
// Returns:
//   - bool: true if the file has an exported date
func (s *State) Exported(filename string) bool {
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
func (s *State) CountUnenriched(journalDir string) int {
	entries, readErr := os.ReadDir(journalDir)
	if readErr != nil {
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
