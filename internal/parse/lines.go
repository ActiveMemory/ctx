//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parse

import "github.com/ActiveMemory/ctx/internal/config/token"

// ByteLines splits data on newline bytes, returning non-empty byte slices.
//
// Parameters:
//   - data: Raw byte content to split
//
// Returns:
//   - [][]byte: Non-empty line slices
func ByteLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == token.NewlineLF[0] {
			if i > start {
				lines = append(lines, data[start:i])
			}
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
