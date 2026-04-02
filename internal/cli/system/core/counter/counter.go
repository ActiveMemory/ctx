//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package counter

import (
	"strconv"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/io"
	ctxLog "github.com/ActiveMemory/ctx/internal/log/warn"
)

// Read reads an integer counter from a file. Returns 0 if the file
// does not exist or cannot be parsed.
//
// Parameters:
//   - path: Absolute path to the counter file
//
// Returns:
//   - int: Counter value, or 0 on error
func Read(path string) int {
	data, readErr := io.SafeReadUserFile(path)
	if readErr != nil {
		return 0
	}
	n, parseErr := strconv.Atoi(strings.TrimSpace(string(data)))
	if parseErr != nil {
		return 0
	}
	return n
}

// Write writes an integer counter to a file.
//
// Parameters:
//   - path: Absolute path to the counter file
//   - n: Counter value to write
func Write(path string, n int) {
	if writeErr := io.SafeWriteFile(
		path, []byte(strconv.Itoa(n)), fs.PermSecret,
	); writeErr != nil {
		ctxLog.Warn(warn.Write, path, writeErr)
	}
}
