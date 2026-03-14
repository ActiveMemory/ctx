//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package token

const (
	// NewlineCRLF is the Windows new line.
	//
	// We check NewlineCRLF first, then NewlineLF to handle both formats.
	NewlineCRLF = "\r\n"
	// NewlineLF is Unix new line.
	NewlineLF = "\n"
	// Whitespace is the set of inline whitespace characters (space and tab).
	Whitespace = " \t"
	// Space is a single space character.
	Space = " "
	// Tab is a horizontal tab character.
	Tab = "\t"
)
