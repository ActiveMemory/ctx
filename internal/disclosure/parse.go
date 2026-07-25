//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure

import (
	cfgDisc "github.com/ActiveMemory/ctx/internal/config/disclosure"
)

// Parse splits a root's content into its regions (preamble, staging,
// themes) for the given kind. It is total: any byte string splits into
// segments, so parsing never fails — all refusal lives in Validate.
//
// The raw segments are kept verbatim, so Reconstruct returns the input
// byte-for-byte; nothing is normalized. Themes is the parsed view of the
// raw themes region. HasThemes is false for a not-yet-migrated root.
//
// Layout is the same for every kind (see specs/progressive-disclosure.md):
//
//	preamble | staging | ## Themes …
//
// Only the prefix opening a staged entry differs — "## [" for the
// timestamped kinds, "## " for a convention's curated prose sections.
//
// Parameters:
//   - content: the full root file content
//   - k: which canonical file this is
//
// Returns:
//   - Root: the split root; Reconstruct(r) == content
func Parse(content string, k Kind) Root {
	r := Root{Kind: k}
	themeOffsets := headingLineOffsets(content, cfgDisc.HeadingThemes)
	r.HasThemes = len(themeOffsets) > 0

	parseRegions(&r, content, themeOffsets)

	r.Themes = parseThemes(r.ThemesRaw)
	return r
}

// Reconstruct returns the root's content in file order. It is the inverse
// of Parse: Reconstruct(Parse(c, k)) == c. The region order is the same
// for every kind, so this needs no per-kind branch.
//
// Returns:
//   - string: the reassembled root content
func (r Root) Reconstruct() string {
	return r.Preamble + r.Staging + r.ThemesRaw
}
