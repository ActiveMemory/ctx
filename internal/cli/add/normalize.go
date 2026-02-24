//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config"
)

// normalizeTargetSection ensures a section heading has a proper Markdown
// format.
//
// Prepends "## " if the section string does not already start with "##".
// Callers must not pass an empty string; the empty case is handled by
// insertTask before this function is reached.
//
// Parameters:
//   - section: Raw section name from user input (non-empty)
//
// Returns:
//   - string: Normalized section heading (e.g., "## Phase 1")
//
// checkRequired returns the names of any fields whose values are empty.
func checkRequired(fields [][2]string) []string {
	var missing []string
	for _, f := range fields {
		if f[1] == "" {
			missing = append(missing, f[0])
		}
	}
	return missing
}

func normalizeTargetSection(section string) string {
	if !strings.HasPrefix(section, config.HeadingLevelTwoStart) {
		return config.HeadingLevelTwoStart + section
	}
	return section
}
