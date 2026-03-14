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

// Truncate returns the first line of s, capped at max characters.
//
// Parameters:
//   - s: input string (may be multi-line).
//   - max: maximum length including ellipsis.
//
// Returns:
//   - string: truncated first line.
func Truncate(s string, max int) string {
	line, _, _ := strings.Cut(s, token.NewlineLF)
	if len(line) <= max {
		return line
	}
	return line[:max-len(token.Ellipsis)] + token.Ellipsis
}
