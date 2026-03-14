//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/token"
)

// ParseEntries splits raw bytes into entry lines, filtering empty lines.
//
// Parameters:
//   - data: Raw scratchpad content
//
// Returns:
//   - []string: Non-empty lines, or nil if data is empty
func ParseEntries(data []byte) []string {
	if len(data) == 0 {
		return nil
	}
	lines := strings.Split(string(data), token.NewlineLF)
	var entries []string
	for _, line := range lines {
		if line != "" {
			entries = append(entries, line)
		}
	}
	return entries
}

// FormatEntries joins entries with newlines and adds a trailing newline.
//
// Parameters:
//   - entries: The scratchpad entries to serialize
//
// Returns:
//   - []byte: Newline-delimited content, or nil if entries is empty
func FormatEntries(entries []string) []byte {
	if len(entries) == 0 {
		return nil
	}
	return []byte(strings.Join(entries, token.NewlineLF) + token.NewlineLF)
}
