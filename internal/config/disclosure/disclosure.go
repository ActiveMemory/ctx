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

	// HeadingRecent delimits the staging zone of a CONVENTIONS root.
	// Conventions append at EOF and their entries are prose sections, so
	// their staging needs an explicit trailing heading that "## Themes"
	// cannot provide.
	HeadingRecent = "## Recent"

	// EntryLinePrefix marks a timestamped entry heading ("## [ts] Title")
	// in LEARNINGS/DECISIONS.
	EntryLinePrefix = "## ["

	// ConventionLinePrefix marks a prose section heading in CONVENTIONS.
	ConventionLinePrefix = "### "

	// IDSeparator joins the timestamp and title of an entry identity. A
	// NUL never appears in a heading line, so entry text cannot forge it.
	IDSeparator = "\x00"

	// LinkOpen is the "](" that separates a markdown link's label from its
	// target; a theme gist's link is the "(...)" following it.
	LinkOpen = "]("

	// ThemeDirLearning, ThemeDirDecision, and ThemeDirConvention name the
	// per-kind subdirectories of the context directory that hold theme
	// files (<theme>.md), reachable only via the root's links.
	ThemeDirLearning   = "learnings"
	ThemeDirDecision   = "decisions"
	ThemeDirConvention = "conventions"
)
