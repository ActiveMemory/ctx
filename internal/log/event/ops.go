//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package event

import (
	"errors"
	"os"

	"github.com/ActiveMemory/ctx/internal/config/event"
)

// rotate checks the current log file size and renames it to the
// previous-generation path when it exceeds [event.LogMaxBytes].
//
// Returns nil when there is nothing to do (log missing or under the
// size threshold), or when rotation succeeds. Any other failure is
// propagated so callers honour the log-first principle: if the log
// cannot be rotated, [Append] cannot safely continue writing and
// downstream side effects must not fire.
//
// Parameters:
//   - logPath: absolute path to the current event log
//
// Returns:
//   - error: [os.ErrNotExist] from the Stat or Remove path is treated
//     as "nothing to rotate" / "nothing to clean up" and returns nil.
//     Any other stat, path, rename, or remove failure is surfaced.
func rotate(logPath string) error {
	info, statErr := os.Stat(logPath)
	if statErr != nil {
		if errors.Is(statErr, os.ErrNotExist) {
			return nil // nothing to rotate yet
		}
		return statErr
	}
	if info.Size() < int64(event.LogMaxBytes) {
		return nil
	}

	prevPath, prevErr := prevLogFilePath()
	if prevErr != nil {
		return prevErr
	}
	if removeErr := os.Remove(prevPath); removeErr != nil {
		if !errors.Is(removeErr, os.ErrNotExist) {
			return removeErr
		}
		// ErrNotExist is fine: no previous generation to remove.
	}
	return os.Rename(logPath, prevPath)
}
