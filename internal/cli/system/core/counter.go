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

	"github.com/ActiveMemory/ctx/internal/io"
)

// ReadCounter reads an integer counter from a file. Returns 0 if the file
// does not exist or cannot be parsed.
//
// Parameters:
//   - path: Absolute path to the counter file
//
// Returns:
//   - int: Counter value, or 0 on error
func ReadCounter(path string) int {
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

// WriteCounter writes an integer counter to a file.
//
// Parameters:
//   - path: Absolute path to the counter file
//   - n: Counter value to write
func WriteCounter(path string, n int) {
	_ = os.WriteFile(path, []byte(strconv.Itoa(n)), 0o600)
}
