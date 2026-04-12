//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// loadState reads sync state from .context/hub/.
// Acquires a lock file to prevent concurrent access.
//
// Returns:
//   - state: Loaded sync state (zero value if no file exists)
//   - func(): Release function to remove the lock file
//   - error: Non-nil on I/O or lock-contention failure
func loadState() (state, func(), error) {
	var s state
	dir := filepath.Join(rc.ContextDir(), cfgHub.DirHub)
	lockPath := filepath.Join(dir, cfgHub.FileSyncLock)

	if mkErr := io.SafeMkdirAll(
		dir, fs.PermKeyDir,
	); mkErr != nil {
		return s, nil, mkErr
	}

	// Acquire lock: fail if another sync is running.
	if _, statErr := os.Stat(lockPath); statErr == nil {
		return s, nil, os.ErrExist
	}
	if writeErr := io.SafeWriteFile(
		lockPath, []byte(cfgHub.LockSentinel), fs.PermFile,
	); writeErr != nil {
		return s, nil, writeErr
	}

	release := func() { _ = os.Remove(lockPath) }

	path := filepath.Join(dir, cfgHub.FileSyncState)
	data, readErr := io.SafeReadUserFile(path)
	if os.IsNotExist(readErr) {
		return s, release, nil
	}
	if readErr != nil {
		release()
		return s, nil, readErr
	}
	if len(data) == 0 {
		return s, release, nil
	}
	if unmarshalErr := json.Unmarshal(
		data, &s,
	); unmarshalErr != nil {
		release()
		return s, nil, unmarshalErr
	}
	return s, release, nil
}

// saveState writes sync state to .context/hub/.
//
// Parameters:
//   - s: State to persist
//
// Returns:
//   - error: Non-nil on marshal or I/O failure
func saveState(s state) error {
	dir := filepath.Join(rc.ContextDir(), cfgHub.DirHub)
	data, marshalErr := json.MarshalIndent(
		s, "", cfgHub.JSONIndent,
	)
	if marshalErr != nil {
		return marshalErr
	}
	path := filepath.Join(dir, cfgHub.FileSyncState)
	return io.SafeWriteFile(path, data, fs.PermFile)
}
