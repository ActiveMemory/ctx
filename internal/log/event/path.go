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
func logFilePath() string {
	return filepath.Join(rc.ContextDir(), dir.State, event.FileLog)
}

// prevLogFilePath returns the absolute path to the rotated event log.
//
// Returns:
//   - string: path under the active context directory
func prevLogFilePath() string {
	return filepath.Join(rc.ContextDir(), dir.State, event.FileLogPrev)
}
