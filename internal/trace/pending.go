//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package trace

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	cfgFs "github.com/ActiveMemory/ctx/internal/config/fs"
	cfgTrace "github.com/ActiveMemory/ctx/internal/config/trace"
)

// Record appends a single pending context reference to the pending file
// in stateDir. The operation is best-effort: directory creation and
// write errors are returned to the caller but do not panic.
//
// Parameters:
//   - ref: context reference string (e.g. "T-3", "D-1", "L-5")
//   - stateDir: absolute path to the state directory
//
// Returns:
//   - error: non-nil if the directory cannot be created or the entry
//     cannot be written
func Record(ref, stateDir string) error {
	entry := PendingEntry{
		Ref:       ref,
		Timestamp: time.Now().UTC(),
	}

	return appendJSONL(stateDir, cfgTrace.FilePending, entry)
}

// ReadPending reads all pending context reference entries from the
// pending file in stateDir. Malformed JSONL lines are silently skipped.
// Returns an empty (non-nil) slice when the file does not exist.
//
// Parameters:
//   - stateDir: absolute path to the state directory
//
// Returns:
//   - []PendingEntry: entries in file order
//   - error: non-nil only if the file exists but cannot be opened
func ReadPending(stateDir string) ([]PendingEntry, error) {
	path := filepath.Join(stateDir, cfgTrace.FilePending)
	return readJSONL[PendingEntry](path)
}

// TruncatePending empties the pending file in stateDir without deleting
// it. If the file does not exist the call is a no-op.
//
// Parameters:
//   - stateDir: absolute path to the state directory
//
// Returns:
//   - error: non-nil if the file exists but cannot be truncated
func TruncatePending(stateDir string) error {
	path := filepath.Join(stateDir, cfgTrace.FilePending)
	//nolint:gosec // path built from trusted stateDir + constant filename
	f, openErr := os.OpenFile(
		filepath.Clean(path),
		os.O_TRUNC|os.O_WRONLY,
		cfgFs.PermFile,
	)
	if openErr != nil {
		if errors.Is(openErr, os.ErrNotExist) {
			return nil
		}
		return openErr
	}
	return f.Close()
}
