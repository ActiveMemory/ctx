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
	"github.com/ActiveMemory/ctx/internal/config/warn"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
	ctxLog "github.com/ActiveMemory/ctx/internal/log/warn"
)

// Message appends a timestamped log line to the given file.
// Rotates the log when it exceeds config.HookLogMaxBytes, keeping one
// previous generation (.1 suffix) - same pattern as eventlog.
//
// Message is a verbose operational logger for hooks: it records what
// a hook did for later debugging, not an authoritative audit trail.
// Unlike [event.Append], a failed write here is not a correctness
// problem for downstream; callers do not gate side effects on the
// line landing. Write errors are therefore logged to stderr via
// [log/warn.Warn] instead of propagated, preserving the void
// signature and the "never break the editor" contract for hooks.
//
// Parameters:
//   - logFile: Absolute path to the log file
//   - sessionID: Session identifier (truncated to 8 chars)
//   - msg: Log message to append
func Message(logFile, sessionID, msg string) {
	d := filepath.Dir(logFile)
	mkdirErr := internalIo.SafeMkdirAll(d, fs.PermRestrictedDir)
	if mkdirErr != nil {
		ctxLog.Warn(warn.Mkdir, d, mkdirErr)
	}

	Rotate(logFile)

	short := sessionID
	if len(short) > journal.SessionIDShortLen {
		short = short[:journal.SessionIDShortLen]
	}

	line := fmt.Sprintf(desc.Text(text.DescKeyWriteLogLineFormat),
		time.Now().Format(cfgTime.DateTimePreciseFmt), short, msg)

	if appendErr := internalIo.AppendBytes(
		logFile, []byte(line), fs.PermSecret,
	); appendErr != nil {
		ctxLog.Warn(warn.Write, logFile, appendErr)
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
		ctxLog.Warn(warn.Remove, prev, removeErr)
	}
	if renameErr := os.Rename(logFile, prev); renameErr != nil {
		ctxLog.Warn(
			warn.Rename, logFile, renameErr,
		)
	}
}
