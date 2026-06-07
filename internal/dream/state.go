//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dream

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"sort"

	cfgDream "github.com/ActiveMemory/ctx/internal/config/dream"
	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	errDream "github.com/ActiveMemory/ctx/internal/err/dream"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// Hash computes the content hash the discipline clock compares against:
// the full SHA-256 hex digest of content. A source re-enters triage
// only when this digest differs from the one saved in state.
//
// Parameters:
//   - content: the raw file bytes to hash
//
// Returns:
//   - string: lowercase hex SHA-256 digest of content
func Hash(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}

// LoadState reads and decodes the per-source state slice from
// <dreamsDir>/state.json. A missing file is not an error: it yields an
// empty slice (the first-ever pass has no saved state).
//
// Parameters:
//   - dreamsDir: absolute path to the dreams/ notebook directory
//
// Returns:
//   - []SourceState: the saved per-source records (empty when none)
//   - error: non-nil on a read failure other than not-exist, or on a
//     JSON decode failure
func LoadState(dreamsDir string) ([]SourceState, error) {
	path := statePath(dreamsDir)
	data, readErr := ctxIo.SafeReadUserFile(path)
	if readErr != nil {
		if os.IsNotExist(readErr) {
			return []SourceState{}, nil
		}
		return nil, errDream.ReadState(path, readErr)
	}
	var states []SourceState
	if unmarshalErr := json.Unmarshal(data, &states); unmarshalErr != nil {
		return nil, errDream.UnmarshalState(path, unmarshalErr)
	}
	return states, nil
}

// SaveState encodes states as indented JSON and writes it atomically to
// <dreamsDir>/state.json, creating the notebook directory if needed.
//
// Parameters:
//   - dreamsDir: absolute path to the dreams/ notebook directory
//   - states: the per-source records to persist
//
// Returns:
//   - error: non-nil on directory creation, JSON marshal, or write
//     failure
func SaveState(dreamsDir string, states []SourceState) error {
	if mkErr := ctxIo.SafeMkdirAll(
		dreamsDir, cfgFs.PermRestrictedDir,
	); mkErr != nil {
		return errDream.Mkdir(dreamsDir, mkErr)
	}
	data, marshalErr := json.MarshalIndent(
		states, "", cfgDream.JSONIndent,
	)
	if marshalErr != nil {
		return errDream.MarshalState(marshalErr)
	}
	path := statePath(dreamsDir)
	if writeErr := ctxIo.SafeWriteFileAtomic(
		path, data, cfgFs.PermSecret,
	); writeErr != nil {
		return errDream.WriteState(path, writeErr)
	}
	return nil
}

// DeltaSelect is the discipline clock: given the current ideas files
// keyed by path to their content hash, it returns the paths whose hash
// is new (no saved record) or changed (hash differs from the saved
// record) versus prior state. Unchanged-and-already-recorded sources
// are skipped. The result is sorted for deterministic ordering.
//
// Parameters:
//   - prior: the previously saved per-source records
//   - current: path → content hash for the ideas files scanned this pass
//
// Returns:
//   - []string: sorted paths that are new or changed since prior state
func DeltaSelect(
	prior []SourceState, current map[string]string,
) []string {
	priorByPath := make(map[string]string, len(prior))
	for _, s := range prior {
		priorByPath[s.Path] = s.Hash
	}
	var selected []string
	for path, hash := range current {
		savedHash, seen := priorByPath[path]
		if !seen || savedHash != hash {
			selected = append(selected, path)
		}
	}
	sort.Strings(selected)
	return selected
}
