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
	ModeReplace Mode = iota
	ModeAppend
	ModePrepend
	ModeBlob
)

// Opts holds all parameters for an edit operation.
type Opts struct {
	N         int
	Text      string
	FilePath  string
	LabelText string
	Mode      Mode
}
