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

// CheckRequired returns the names of any fields whose values are empty.
//
// Parameters:
//   - fields: Pairs of [name, value] to check
//
// Returns:
//   - []string: Names of fields with empty values
func CheckRequired(fields [][2]string) []string {
	var missing []string
	for _, f := range fields {
		if f[1] == "" {
			missing = append(missing, f[0])
		}
	}
	return missing
}

// NormalizeTargetSection ensures a section heading has proper Markdown format.
//
// Prepends "## " if the section string does not already start with "##".
// Callers must not pass an empty string; the empty case is handled by
// InsertTask before this function is reached.
//
// Parameters:
//   - section: Raw section name from user input (non-empty)
//
// Returns:
//   - string: Normalized section heading (e.g., "## Phase 1")
func NormalizeTargetSection(section string) string {
	if !strings.HasPrefix(section, token.HeadingLevelTwoStart) {
		return token.HeadingLevelTwoStart + section
	}
	return section
}
