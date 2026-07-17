//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure

// Kind identifies which canonical knowledge file a root is, because the
// staging zone sits in a different place per kind (above ## Themes for
// entry files, inside ## Recent for conventions).
type Kind int

const (
	// KindLearning is LEARNINGS.md — "## [ts] Title" entries, staging
	// above ## Themes.
	KindLearning Kind = iota
	// KindDecision is DECISIONS.md — same entry shape as learnings.
	KindDecision
	// KindConvention is CONVENTIONS.md — prose sections, staging in
	// ## Recent.
	KindConvention
)

// Theme is one entry in a root's ## Themes section: a name, the
// "just enough" gist, and the markdown link to the theme file holding
// that theme's bodies.
//
// Fields:
//   - Name: the theme's short name (heading text)
//   - Gist: the "just enough" one-liner conveying the theme
//   - Link: relative path to the theme file holding the bodies
type Theme struct {
	Name string
	Gist string
	Link string
}

// Root is a parsed progressive-disclosure root, split into its regions.
//
// The raw segments (Preamble, Staging, ThemesRaw) are kept verbatim so a
// root round-trips byte-for-byte via Reconstruct: parsing never
// normalizes content, which is how the clobber-bug class is avoided.
// Themes is the parsed view of ThemesRaw, for consumers that read the
// gists.
//
// HasThemes is false for a not-yet-migrated root (no ## Themes section);
// Validate treats that as the first-run case, not an error.
//
// Fields:
//   - Kind: which canonical file this root is
//   - Preamble: everything above the first region (H1 + comment block)
//   - Staging: the un-digested entries region (verbatim)
//   - ThemesRaw: the raw ## Themes region (verbatim), empty if none
//   - Themes: parsed view of ThemesRaw
//   - HasThemes: false when the root is not yet migrated
type Root struct {
	Kind      Kind
	Preamble  string
	Staging   string
	ThemesRaw string
	Themes    []Theme
	HasThemes bool
}
