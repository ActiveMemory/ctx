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
	"github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// stateFile is the sync state persistence filename.
const stateFile = ".sync-state.json"

// lockFile is the lock file to prevent concurrent syncs.
const lockFile = ".sync.lock"

// loadState reads sync state from .context/shared/.
// Acquires a lock file to prevent concurrent access.
func loadState() (state, func(), error) {
	var s state
	dir := filepath.Join(rc.ContextDir(), "shared")
	lockPath := filepath.Join(dir, lockFile)

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
		lockPath, []byte("lock"), fs.PermFile,
	); writeErr != nil {
		return s, nil, writeErr
	}

	release := func() { _ = os.Remove(lockPath) }

	path := filepath.Join(dir, stateFile)
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

// saveState writes sync state to .context/shared/.
func saveState(s state) error {
	dir := filepath.Join(rc.ContextDir(), "shared")
	data, marshalErr := json.MarshalIndent(
		s, "", "  ",
	)
	if marshalErr != nil {
		return marshalErr
	}
	path := filepath.Join(dir, stateFile)
	return io.SafeWriteFile(path, data, fs.PermFile)
}
