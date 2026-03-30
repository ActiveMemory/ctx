//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package box

// Nudge box drawing constants.
const (
	// Top is the top-left corner of a nudge box.
	Top = "┌─ "
	// LinePrefix is the left border prefix for nudge box content lines.
	LinePrefix = "│ "
	// Bottom is the bottom border of a nudge box.
	Bottom = "└──────────────────────────────────────────────────"
	// NudgeBoxWidth is the inner character width of the nudge box border.
	NudgeBoxWidth = 51
	// BorderFill is the repeating character used to pad the top border.
	BorderFill = "─"
)

// PipeSeparator is the inline separator used between navigation links.
const PipeSeparator = " | "
