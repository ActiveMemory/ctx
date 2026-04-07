//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parser

import (
	"bufio"
	"os"

	cfgCopilot "github.com/ActiveMemory/ctx/internal/config/copilot"
	"github.com/ActiveMemory/ctx/internal/io"
)

// openScanner opens a JSONL file and returns a buffered
// scanner. The caller must close the returned file.
//
// Parameters:
//   - path: filesystem path to the JSONL file
//   - bufMax: maximum token size for the scanner buffer
//
// Returns:
//   - *os.File: opened file handle (caller must close)
//   - *bufio.Scanner: scanner configured with bufMax
//   - error: file-open failure
func openScanner(
	path string, bufMax int,
) (*os.File, *bufio.Scanner, error) {
	f, openErr := io.SafeOpenUserFile(path)
	if openErr != nil {
		return nil, nil, openErr
	}
	scanner := bufio.NewScanner(f)
	buf := make([]byte, 0, cfgCopilot.ScanBufInit)
	scanner.Buffer(buf, bufMax)
	return f, scanner, nil
}
