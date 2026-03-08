//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"os"
	"strconv"
	"strings"
)

// ReadMtime reads a stored mtime value from a file. Returns 0 if the
// file does not exist or the content cannot be parsed as an integer.
//
// Parameters:
//   - path: absolute path to the mtime state file
//
// Returns:
//   - int64: the stored mtime value, or 0 on error
func ReadMtime(path string) int64 {
	data, readErr := os.ReadFile(path) //nolint:gosec // temp file path
	if readErr != nil {
		return 0
	}
	n, parseErr := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
	if parseErr != nil {
		return 0
	}
	return n
}

// WriteMtime writes a mtime value to a file. Errors are silently
// ignored (best-effort state persistence).
//
// Parameters:
//   - path: absolute path to the mtime state file
//   - mtime: the mtime value to store
func WriteMtime(path string, mtime int64) {
	_ = os.WriteFile(path, []byte(strconv.FormatInt(mtime, 10)), 0o600)
}
