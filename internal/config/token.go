//   /    Context:                     https://ctx.ist
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
	// Separator is a Markdown horizontal rule used between sections.
	Separator = "---"
	// Ellipsis is a Markdown ellipsis.
	Ellipsis = "..."
	// HeadingLevelOneStart is the Markdown heading for the first section.
	HeadingLevelOneStart = "# "
	// HeadingLevelTwoStart is the Markdown heading for subsequent sections.
	HeadingLevelTwoStart = "## "
)
