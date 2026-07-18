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
// Layout per kind (see specs/progressive-disclosure.md):
//   - entry kinds:  preamble | staging (## [ entries) | ## Themes …
//   - conventions:  preamble | ## Themes … | ## Recent (staging)
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

	if k == KindConvention {
		parseConvention(&r, content, themeOffsets)
	} else {
		parseEntryKind(&r, content, themeOffsets)
	}

	r.Themes = parseThemes(r.ThemesRaw)
	return r
}

// Reconstruct returns the root's content in file order for its kind. It
// is the inverse of Parse: Reconstruct(Parse(c, k)) == c.
//
// Returns:
//   - string: the reassembled root content
func (r Root) Reconstruct() string {
	if r.Kind == KindConvention {
		return r.Preamble + r.ThemesRaw + r.Staging
	}
	return r.Preamble + r.Staging + r.ThemesRaw
}
