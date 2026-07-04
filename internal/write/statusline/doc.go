//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package statusline emits the rendered status line for the hidden
// "ctx system statusline" command.
//
// Claude Code displays the first stdout line of the configured
// statusLine command; these primitives are the single place that
// line is written.
//
// # Exported Functions
//
// [Line] prints the assembled status line.
//
// [Blank] prints an empty line, used when statusline.enabled is
// false so the displayed status line goes blank without a non-zero
// exit.
//
// # Message Categories
//
//   - Info: the rendered status line
package statusline
