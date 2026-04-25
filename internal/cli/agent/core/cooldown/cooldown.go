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

	"github.com/ActiveMemory/ctx/internal/config/agent"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/rc"
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
//   - error: non-nil when the context directory is not declared or
//     the state directory cannot be created. Previously this helper
//     logged the mkdir error and returned the path anyway, guaranteeing
//     a second failure on the subsequent write; propagating keeps the
//     first failure authoritative.
func TombstonePath(session string) (string, error) {
	ctxDir, err := rc.ContextDir()
	if err != nil {
		return "", err
	}
	stateDir := filepath.Join(ctxDir, dir.State)
	if mkdirErr := ctxIo.SafeMkdirAll(
		stateDir, fs.PermRestrictedDir,
	); mkdirErr != nil {
		return "", mkdirErr
	}
	return filepath.Join(stateDir, agent.TombstonePrefix+session), nil
}
