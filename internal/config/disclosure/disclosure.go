//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package disclosure

// Structural vocabulary for progressive-disclosure roots: the headings
// that delimit a root's regions, the line prefixes that mark entries,
// and the per-kind theme-file directories.
const (
	// HeadingThemes delimits the themes region of a root — the per-theme
	// gists and links. In entry-based roots (LEARNINGS/DECISIONS) it also
	// marks the lower bound of the staging zone: no entry may appear below
	// it.
	HeadingThemes = "## Themes"

	// EntryLinePrefix marks a timestamped entry heading ("## [ts] Title")
	// in LEARNINGS/DECISIONS.
	EntryLinePrefix = "## ["

	// SectionLinePrefix marks a curated prose section heading ("## Title")
	// in CONVENTIONS — the convention kind's entry prefix, the counterpart
	// to EntryLinePrefix for timestamped kinds. HeadingThemes shares this
	// prefix, so callers scanning for entries must skip it explicitly.
	SectionLinePrefix = "## "

	// SectionUnfiled names the section used only when a non-CLI caller
	// reaches the convention add-path with no section and the root holds
	// none either. It is NOT a default the CLI offers: `ctx convention
	// add` refuses an empty or placeholder --section outright, because a
	// catch-all section is where an undecided caller dumps everything.
	// The name is deliberately unattractive so a human notices it.
	SectionUnfiled = "Unfiled"

	// IDSeparator joins the timestamp and title of an entry identity. A
	// NUL never appears in a heading line, so entry text cannot forge it.
	IDSeparator = "\x00"

	// LinkOpen is the "](" that separates a markdown link's label from its
	// target; a theme gist's link is the "(...)" following it.
	LinkOpen = "]("

	// LinkLabelOpen is the "[" that begins a markdown link's label.
	LinkLabelOpen = "["

	// LinkClose is the ")" that ends a markdown link's target.
	LinkClose = ")"

	// ThemeArrow separates a theme bullet's gist from its markdown link:
	// "- name — gist → [name](noun/slug.md)". The mover's gist write-back
	// renders it; parseThemeBullet locates the link via LinkOpen, so the
	// arrow is cosmetic structure rather than a parse anchor.
	ThemeArrow = " → "

	// ThemeDirLearning, ThemeDirDecision, and ThemeDirConvention name the
	// per-kind subdirectories of the context directory that hold theme
	// files (<theme>.md), reachable only via the root's links.
	ThemeDirLearning   = "learnings"
	ThemeDirDecision   = "decisions"
	ThemeDirConvention = "conventions"
)
