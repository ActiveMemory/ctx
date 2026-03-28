//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package format

import "github.com/ActiveMemory/ctx/internal/config/token"

// Truncate shortens s to max characters, appending an ellipsis if truncated.
//
// Parameters:
//   - s: String to truncate
//   - max: Maximum length including ellipsis
//
// Returns:
//   - string: Truncated string
func Truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-len(token.Ellipsis)] + token.Ellipsis
}
