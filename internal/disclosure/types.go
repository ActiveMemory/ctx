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
	Name string `json:"name"`
	Gist string `json:"gist"`
	Link string `json:"link"`
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

// StagedEntry is one un-digested entry in a root's staging zone,
// identified for the dry-run pass to propose a theme for.
//
// Fields:
//   - Timestamp: the entry's timestamp (e.g. "2026-07-18-120000")
//   - Title: the entry's title text
type StagedEntry struct {
	Timestamp string `json:"timestamp"`
	Title     string `json:"title"`
}

// Inspection is the read-only view of a root the dry-run pass consumes:
// what kind it is, which entries are staged (awaiting digestion), and
// which themes already exist. It is the structured form of a Parse,
// stable across the CLI's JSON output.
//
// Fields:
//   - Kind: the root's kind name ("learning" | "decision" | "convention")
//   - Staging: the un-digested entries, in file order
//   - Themes: the current ## Themes (name/gist/link), in file order
type Inspection struct {
	Kind    string        `json:"kind"`
	Staging []StagedEntry `json:"staging"`
	Themes  []Theme       `json:"themes"`
}

// Plan is the digest plan the ctx-digest skill authors and Apply
// executes: per target theme, which staged entries move there and the
// gist to write back. Entry identity is timestamp+title (joined by
// cfgDisc.IDSeparator), matching entryIDs and CheckUniqueness.
//
// Fields:
//   - Kind: the root's kind name ("learning" | "decision")
//   - Assignments: one per target theme; together they partition the
//     entries the pass moves out of staging
type Plan struct {
	Kind        string       `json:"kind"`
	Assignments []Assignment `json:"assignments"`
}

// Assignment moves a set of staged entries into one theme and (re)writes
// that theme's gist bullet in the root's ## Themes. Its Entries are the
// same {timestamp, title} shape inspect reports under "staging", so a
// digest plan lifts them verbatim from the inspection — no NUL-joined id
// to hand-author.
//
// Fields:
//   - Theme: the theme's name — bullet label, and heading on first use
//   - Slug: the theme-file basename stem, resolved to <noun>/<slug>.md
//   - Gist: the authored one-line gist (spec ### Gist format)
//   - Entries: the staged entries to move here, in file order
type Assignment struct {
	Theme   string        `json:"theme"`
	Slug    string        `json:"slug"`
	Gist    string        `json:"gist"`
	Entries []StagedEntry `json:"entries"`
}

// ApplyResult reports what a successful Apply did, for CLI output.
//
// Fields:
//   - Moved: the number of entries moved out of staging
//   - Themes: the theme slugs created or appended, in plan order
type ApplyResult struct {
	Moved  int      `json:"moved"`
	Themes []string `json:"themes"`
}
