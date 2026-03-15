//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package log

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/ActiveMemory/ctx/internal/config/event"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// logFilePath returns the absolute path to the current event log
// under the active context directory.
func logFilePath() string {
	return filepath.Join(rc.ContextDir(), dir.State, event.FileEventLog)
}

// prevLogFilePath returns the absolute path to the rotated (previous
// generation) event log under the active context directory.
func prevLogFilePath() string {
	return filepath.Join(rc.ContextDir(), dir.State, event.FileEventLogPrev)
}
