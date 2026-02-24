//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package context

// splitLines splits content into lines without using strings package.
//
// Parameters:
//   - content: Byte slice to split
//
// Returns:
//   - [][]byte: Slice of lines (without newline characters)
func splitLines(content []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range content {
		if b == '\n' {
			lines = append(lines, content[start:i])
			start = i + 1
		}
	}
	if start < len(content) {
		lines = append(lines, content[start:])
	}
	return lines
}

// trimSpace trims leading and trailing whitespace from a byte slice.
//
// Parameters:
//   - b: Byte slice to trim
//
// Returns:
//   - []byte: Slice with spaces, tabs, and carriage returns removed from ends
func trimSpace(b []byte) []byte {
	start := 0
	end := len(b)
	for start < end && (b[start] == ' ' || b[start] == '\t' || b[start] == '\r') {
		start++
	}
	for end > start && (b[end-1] == ' ' || b[end-1] == '\t' || b[end-1] == '\r') {
		end--
	}
	return b[start:end]
}
