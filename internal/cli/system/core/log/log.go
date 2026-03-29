//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package log

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/event"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/journal"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	ctxLog "github.com/ActiveMemory/ctx/internal/log"
)

// Message appends a timestamped log line to the given file.
// Rotates the log when it exceeds config.HookLogMaxBytes, keeping one
// previous generation (.1 suffix) - same pattern as eventlog.
//
// Parameters:
//   - logFile: Absolute path to the log file
//   - sessionID: Session identifier (truncated to 8 chars)
//   - msg: Log message to append
func Message(logFile, sessionID, msg string) {
	d := filepath.Dir(logFile)
	if mkdirErr := os.MkdirAll(d, fs.PermRestrictedDir); mkdirErr != nil {
		ctxLog.Warn("mkdir %s: %v", d, mkdirErr)
	}

	Rotate(logFile)

	short := sessionID
	if len(short) > journal.SessionIDShortLen {
		short = short[:journal.SessionIDShortLen]
	}

	line := fmt.Sprintf(desc.Text(text.DescKeyWriteLogLineFormat),
		time.Now().Format(cfgTime.DateTimePreciseFmt), short, msg)

	f, openErr := os.OpenFile( //nolint:gosec // logFile is constructed internally
		logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		fs.PermSecret,
	)
	if openErr != nil {
		return
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			ctxLog.Warn("close %s: %v", logFile, closeErr)
		}
	}()
	if _, writeErr := f.WriteString(line); writeErr != nil {
		ctxLog.Warn("write %s: %v", logFile, writeErr)
	}
}

// Rotate checks the log file size and rotates if it exceeds
// config.HookLogMaxBytes. The previous generation is replaced.
//
// Parameters:
//   - logFile: Absolute path to the log file
func Rotate(logFile string) {
	info, statErr := os.Stat(logFile)
	if statErr != nil {
		return
	}
	if info.Size() < int64(event.HookLogMaxBytes) {
		return
	}
	prev := logFile + event.RotationSuffix
	if removeErr := os.Remove(prev); removeErr != nil {
		ctxLog.Warn("remove %s: %v", prev, removeErr)
	}
	if renameErr := os.Rename(logFile, prev); renameErr != nil {
		ctxLog.Warn(
			"rename %s: %v", logFile, renameErr,
		)
	}
}
