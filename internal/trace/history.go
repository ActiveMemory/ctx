//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"path/filepath"
	"strings"
	"time"
)

// historyFile is the name of the JSONL file that stores commit history entries.
const historyFile = "history.jsonl"

// overrideFile is the name of the JSONL file that stores override entries.
const overrideFile = "overrides.jsonl"

// WriteHistory appends a HistoryEntry to history.jsonl in traceDir.
// If entry.Timestamp is zero it is set to the current UTC time.
// The traceDir is created with MkdirAll if it does not exist.
//
// Parameters:
//   - entry: the HistoryEntry to persist
//   - traceDir: absolute path to the trace directory
//
// Returns:
//   - error: non-nil if the directory cannot be created or the entry cannot be written
func WriteHistory(entry HistoryEntry, traceDir string) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}

	return appendJSONL(traceDir, historyFile, entry)
}

// ReadHistory reads all HistoryEntry records from history.jsonl in traceDir.
// Malformed JSONL lines are silently skipped.
// Returns an empty (non-nil) slice when the file does not exist.
//
// Parameters:
//   - traceDir: absolute path to the trace directory
//
// Returns:
//   - []HistoryEntry: entries in file order
//   - error: non-nil only if the file exists but cannot be read
func ReadHistory(traceDir string) ([]HistoryEntry, error) {
	path := filepath.Join(traceDir, historyFile)
	return readJSONL[HistoryEntry](path)
}

// ReadHistoryForCommit finds the first HistoryEntry whose Commit field matches
// commitHash. Matching supports short hashes by checking whether either string
// is a prefix of the other.
//
// Parameters:
//   - commitHash: full or abbreviated commit hash to look up
//   - traceDir: absolute path to the trace directory
//
// Returns:
//   - HistoryEntry: the matching entry (zero value if not found)
//   - bool: true if a match was found
func ReadHistoryForCommit(commitHash, traceDir string) (HistoryEntry, bool) {
	entries, err := ReadHistory(traceDir)
	if err != nil {
		return HistoryEntry{}, false
	}

	for _, e := range entries {
		if matchesCommit(e.Commit, commitHash) {
			return e, true
		}
	}

	return HistoryEntry{}, false
}

// matchesCommit checks whether stored and query commit hashes match.
// Supports short hashes by checking whether either string is a prefix
// of the other.
func matchesCommit(stored, query string) bool {
	return strings.HasPrefix(stored, query) || strings.HasPrefix(query, stored)
}

// WriteOverride appends an OverrideEntry to overrides.jsonl in traceDir.
// If entry.Timestamp is zero it is set to the current UTC time.
// The traceDir is created with MkdirAll if it does not exist.
//
// Parameters:
//   - entry: the OverrideEntry to persist
//   - traceDir: absolute path to the trace directory
//
// Returns:
//   - error: non-nil if the directory cannot be created or the entry cannot be written
func WriteOverride(entry OverrideEntry, traceDir string) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}

	return appendJSONL(traceDir, overrideFile, entry)
}

// ReadOverrides reads all OverrideEntry records from overrides.jsonl in traceDir.
// Malformed JSONL lines are silently skipped.
// Returns an empty (non-nil) slice when the file does not exist.
//
// Parameters:
//   - traceDir: absolute path to the trace directory
//
// Returns:
//   - []OverrideEntry: entries in file order
//   - error: non-nil only if the file exists but cannot be read
func ReadOverrides(traceDir string) ([]OverrideEntry, error) {
	path := filepath.Join(traceDir, overrideFile)
	return readJSONL[OverrideEntry](path)
}

// ReadOverridesForCommit collects all Refs from OverrideEntry records whose
// Commit field matches commitHash. Matching supports short hashes by checking
// whether either string is a prefix of the other. Refs from all matching
// entries are concatenated and returned as a flat list.
//
// Parameters:
//   - commitHash: full or abbreviated commit hash to look up
//   - traceDir: absolute path to the trace directory
//
// Returns:
//   - []string: flattened list of refs from all matching override entries
func ReadOverridesForCommit(commitHash, traceDir string) []string {
	entries, err := ReadOverrides(traceDir)
	if err != nil {
		return []string{}
	}

	var refs []string
	for _, e := range entries {
		if matchesCommit(e.Commit, commitHash) {
			refs = append(refs, e.Refs...)
		}
	}

	if refs == nil {
		return []string{}
	}

	return refs
}
