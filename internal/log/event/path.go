//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package event

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/event"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// logFilePath returns the absolute path to the current event log.
//
// Returns:
//   - string: path under the active context directory
//   - error: non-nil when the context directory is not declared
func logFilePath() (string, error) {
	ctxDir, err := rc.ContextDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(ctxDir, dir.State, event.FileLog), nil
}

// prevLogFilePath returns the absolute path to the rotated event log.
//
// Returns:
//   - string: path under the active context directory
//   - error: non-nil when the context directory is not declared
func prevLogFilePath() (string, error) {
	ctxDir, err := rc.ContextDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(ctxDir, dir.State, event.FileLogPrev), nil
}
