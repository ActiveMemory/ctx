//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

// Scratchpad file constants for .context/ directory.
const (
	// Enc is the encrypted scratchpad file.
	Enc = "scratchpad.enc"
	// EncOurs is the merge conflict "ours" variant.
	EncOurs = Enc + ".ours"
	// EncTheirs is the merge conflict "theirs" variant.
	EncTheirs = Enc + ".theirs"
	// Md is the plaintext scratchpad file.
	Md = "scratchpad.md"
)

// FmtPadEntryID is the format string for rendering a stable
// entry ID prefix, e.g. "[3] some content".
const FmtPadEntryID = "[%d] %s"

// Merge conflict side labels.
const (
	SideOurs   = "OURS"
	SideTheirs = "THEIRS"
)
