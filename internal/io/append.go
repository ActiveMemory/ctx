//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package io

import (
	"os"

	cfgWarn "github.com/ActiveMemory/ctx/internal/config/warn"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
)

// AppendBytes opens path in append mode, writes data, and closes.
// Errors are logged to stderr via log/warn — this is a best-effort
// operation for JSONL event logs and session stats where failures
// should not interrupt the caller.
//
// Parameters:
//   - path: file path to append to (created if missing)
//   - data: bytes to append
//   - perm: file permission bits for creation
func AppendBytes(path string, data []byte, perm os.FileMode) {
	f, openErr := SafeAppendFile(path, perm)
	if openErr != nil {
		return
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			logWarn.Warn(cfgWarn.Close, path, closeErr)
		}
	}()
	if _, writeErr := f.Write(data); writeErr != nil {
		logWarn.Warn(cfgWarn.Write, path, writeErr)
	}
}
