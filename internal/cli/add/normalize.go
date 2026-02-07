//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import "strings"

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
func normalizeTargetSection(section string) string {
	if !strings.HasPrefix(section, "##") {
		return "## " + section
	}
	return section
}
