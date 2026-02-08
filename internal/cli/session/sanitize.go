//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
)

// cleanInsight cleans and truncates an extracted insight.
//
// Removes leading/trailing whitespace and punctuation, and truncates long
// insights to InsightMaxLen characters at a word boundary when possible.
//
// Parameters:
//   - s: Raw insight string to clean
//
// Returns:
//   - string: Cleaned and potentially truncated insight
func cleanInsight(s string) string {
	s = strings.TrimSpace(s)
	// Remove trailing punctuation fragments
	s = strings.TrimRight(s, ".,;:!?")
	// Truncate if too long
	if len(s) > config.InsightMaxLen {
		// Try to cut at word boundary
		idx := strings.LastIndex(s[:config.InsightMaxLen], " ")
		if idx > config.InsightWordBoundaryMin {
			s = s[:idx] + config.Ellipsis
		} else {
			s = s[:config.InsightMaxLen-len(config.Ellipsis)] + config.Ellipsis
		}
	}
	return s
}

// truncate shortens a string to maxLen characters, adding an ellipsis
// if truncated.
//
// Parameters:
//   - s: String to truncate
//   - maxLen: Maximum length including the ellipsis
//
// Returns:
//   - string: Original string if within maxLen, otherwise truncated
//     with ellipsis
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-len(config.Ellipsis)] + config.Ellipsis
}
