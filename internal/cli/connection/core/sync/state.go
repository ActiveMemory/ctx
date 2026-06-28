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
	cfgWarn "github.com/ActiveMemory/ctx/internal/config/warn"
	errHub "github.com/ActiveMemory/ctx/internal/err/hub"
	"github.com/ActiveMemory/ctx/internal/io"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
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
	ctxDir, ctxErr := rc.ContextDir()
	if ctxErr != nil {
		return s, nil, ctxErr
	}
	dir := filepath.Join(ctxDir, cfgHub.DirHub)
	lockPath := filepath.Join(dir, cfgHub.FileSyncLock)

	if mkErr := io.SafeMkdirAll(
		dir, fs.PermKeyDir,
	); mkErr != nil {
		return s, nil, mkErr
	}

	// Acquire lock: fail if another sync is running. The
	// O_CREATE|O_EXCL create-or-fail is a single syscall, so
	// there is no check-to-write window for a concurrent sync
	// to slip through.
	acquired, lockErr := io.SafeTryLock(lockPath, fs.PermFile)
	if lockErr != nil {
		return s, nil, lockErr
	}
	if !acquired {
		// Another sync holds the lock — or a crashed sync left it
		// stale. Name the path so the wedge is self-documenting;
		// the error still wraps os.ErrExist for errors.Is callers.
		return s, nil, errHub.ConnectSyncLocked(lockPath)
	}

	release := func() {
		if rmErr := io.SafeUnlock(lockPath); rmErr != nil {
			logWarn.Warn(cfgWarn.Remove, lockPath, rmErr)
		}
	}

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
	ctxDir, ctxErr := rc.ContextDir()
	if ctxErr != nil {
		return ctxErr
	}
	dir := filepath.Join(ctxDir, cfgHub.DirHub)
	data, marshalErr := json.MarshalIndent(
		s, "", cfgHub.JSONIndent,
	)
	if marshalErr != nil {
		return marshalErr
	}
	path := filepath.Join(dir, cfgHub.FileSyncState)
	return io.SafeWriteFile(path, data, fs.PermFile)
}
