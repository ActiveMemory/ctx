//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package edit

// Mode selects the edit operation.
type Mode int

// Edit modes.
const (
	// ModeReplace overwrites an existing entry.
	ModeReplace Mode = iota
	// ModeAppend adds text after the existing entry.
	ModeAppend
	// ModePrepend adds text before the existing entry.
	ModePrepend
	// ModeBlob imports a file as a blob entry.
	ModeBlob
)

// Opts holds all parameters for an edit operation.
//
// Fields:
//   - N: Entry index (1-based)
//   - Text: New entry text
//   - FilePath: Path for blob import
//   - LabelText: Display label for blob entries
//   - Mode: Edit mode (replace, append, prepend, blob)
type Opts struct {
	N         int
	Text      string
	FilePath  string
	LabelText string
	Mode      Mode
}
