//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package log

import (
	"os"

	"github.com/ActiveMemory/ctx/internal/config/event"
)

// rotate checks the current log file size and renames it to the
// previous-generation path when it exceeds [event.LogMaxBytes].
// Best-effort: all errors are silently ignored so rotation never
// blocks event logging.
//
// Parameters:
//   - logPath: absolute path to the current event log
func rotate(logPath string) {
	info, statErr := os.Stat(logPath)
	if statErr != nil {
		return // file doesn't exist yet, nothing to rotate
	}
	if info.Size() < int64(event.LogMaxBytes) {
		return
	}

	prevPath := prevLogFilePath()
	_ = os.Remove(prevPath)
	_ = os.Rename(logPath, prevPath)
}
