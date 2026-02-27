//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.

package config

const (
	// NewlineCRLF is the Windows new line.
	//
	// We check NewlineCRLF first, then NewlineLF to handle both formats.
	NewlineCRLF = "\r\n"
	// NewlineLF is Unix new line.
	NewlineLF = "\n"
	// Whitespace is the set of inline whitespace characters (space and tab).
	Whitespace = " \t"
	// Separator is a Markdown horizontal rule used between sections.
	Separator = "---"
	// Ellipsis is a Markdown ellipsis.
	Ellipsis = "..."
	// HeadingLevelOneStart is the Markdown heading for the first section.
	HeadingLevelOneStart = "# "
	// HeadingLevelTwoStart is the Markdown heading for subsequent sections.
	HeadingLevelTwoStart = "## "
	// CodeFence is the standard Markdown code fence delimiter.
	CodeFence = "```"
	// Backtick is a single backtick character.
	Backtick = "`"
	// PipeSeparator is the inline separator used between navigation links.
	PipeSeparator = " | "
	// LinkPrefixParent is the relative link prefix to the parent directory.
	LinkPrefixParent = "../"
	// PrefixHeading is the Markdown heading character used for prefix checks.
	PrefixHeading = "#"
	// PrefixBracket is the opening bracket used for placeholder checks.
	PrefixBracket = "["
	// LoopComplete is the banner printed when the loop finishes.
	LoopComplete = "=== Loop Complete ==="
	// TomlNavOpen is the opening bracket for the TOML nav array.
	TomlNavOpen = "nav = ["
	// TomlNavSectionClose closes a nav section group.
	TomlNavSectionClose = "  ]}"
	// TomlNavClose closes the top-level nav array.
	TomlNavClose = "]"
	// NudgeBoxBottom is the bottom border of a nudge/notification box.
	NudgeBoxBottom = "└──────────────────────────────────────────────────"
)
