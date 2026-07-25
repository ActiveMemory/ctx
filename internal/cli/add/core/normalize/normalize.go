//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package normalize

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/token"
)

// TargetSection ensures a section heading has proper Markdown format.
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
func TargetSection(section string) string {
	if !strings.HasPrefix(section, token.HeadingLevelThreeStart) {
		return token.HeadingLevelThreeStart + section
	}
	return section
}

// ConventionSection ensures a convention section name is an H2 heading.
//
// Conventions group bullets under "## " sections, and the digesting pass
// addresses those sections by heading (specs/progressive-disclosure.md),
// so the heading level is load-bearing rather than cosmetic: an "### "
// section would not be seen as a staged entry.
//
// Parameters:
//   - section: Raw section name from user input (non-empty)
//
// Returns:
//   - string: Normalized H2 heading (e.g., "## Go style")
func ConventionSection(section string) string {
	if !strings.HasPrefix(section, token.HeadingLevelTwoStart) {
		return token.HeadingLevelTwoStart + section
	}
	return section
}
