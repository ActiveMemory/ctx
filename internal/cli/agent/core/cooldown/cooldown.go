//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cooldown

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/config/agent"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
)

// Active checks whether the cooldown tombstone for the given
// session is still fresh.
//
// Parameters:
//   - session: session identifier (typically the caller's PID)
//   - cooldown: duration to suppress repeated output
//
// Returns:
//   - bool: true when the tombstone exists and is within the cooldown
//     window. Always false when cooldown is disabled for this call
//     (empty session or non-positive cooldown) or when no tombstone
//     has ever been written.
//   - error: [os.ErrNotExist] is treated as a legitimate "not active"
//     exit condition and NOT returned. Any other failure (context
//     directory undeclared, permission denied, I/O failure) is
//     surfaced so callers do not silently treat it as "not active"
//     and emit output they meant to suppress.
func Active(session string, cooldown time.Duration) (bool, error) {
	if session == "" || cooldown <= 0 {
		return false, nil
	}
	path, pathErr := TombstonePath(session)
	if pathErr != nil {
		return false, pathErr
	}
	info, statErr := os.Stat(path)
	if statErr != nil {
		if errors.Is(statErr, os.ErrNotExist) {
			// No prior emission; legitimately not active.
			return false, nil
		}
		// Permission denied, I/O failure, etc.: surface.
		return false, statErr
	}
	return time.Since(info.ModTime()) < cooldown, nil
}

// TouchTombstone creates or updates the tombstone file for the given
// session, marking the current time as the last emission.
//
// Parameters:
//   - session: session identifier (typically the caller's PID)
//
// Returns:
//   - error: nil on an empty session (no-op). Non-nil when the
//     tombstone path cannot be resolved or the file cannot be
//     written. Callers decide whether a persistence failure
//     warrants aborting the command; this helper no longer
//     logs and swallows on its own.
func TouchTombstone(session string) error {
	if session == "" {
		return nil
	}
	p, pathErr := TombstonePath(session)
	if pathErr != nil {
		return pathErr
	}
	return ctxIo.SafeWriteFile(p, nil, fs.PermSecret)
}

// TombstonePath returns the filesystem path for a session's tombstone.
//
// Parameters:
//   - session: session identifier
//
// Returns:
//   - string: absolute path under the context state directory.
//   - error: non-nil when the context directory is not declared,
//     the project is not initialized, or the state directory cannot
//     be created. Delegates to [state.Dir] so the
//     [errCtx.ErrNotInitialized] gate applies here too (see
//     specs/state-dir-no-mkdir-when-uninitialized.md). Without this,
//     a hook-driven `ctx agent` invocation in a non-ctx project
//     (the cross-IDE Cursor leak path) would mkdir `.context/state/`
//     directly here, bypassing [state.Dir]'s gate.
func TombstonePath(session string) (string, error) {
	stateDir, dirErr := state.Dir()
	if dirErr != nil {
		return "", dirErr
	}
	return filepath.Join(stateDir, agent.TombstonePrefix+session), nil
}
