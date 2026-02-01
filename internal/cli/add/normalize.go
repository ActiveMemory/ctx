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
// If the section is empty, defaults to "## Next Up". If provided without
// a heading prefix, prepends "## " to make it a level-2 heading.
//
// Parameters:
//   - section: Raw section name from user input
//
// Returns:
//   - string: Normalized section heading (e.g., "## Phase 1")
func normalizeTargetSection(section string) string {
	targetSection := section
	if targetSection == "" {
		return "## Next Up"
	}
	if !strings.HasPrefix(targetSection, "##") {
		return "## " + targetSection
	}
	return targetSection
}
