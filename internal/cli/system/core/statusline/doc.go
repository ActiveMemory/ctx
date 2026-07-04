//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package statusline renders the segments of the ctx status line
// from Claude Code's statusLine stdin payload.
//
// The cmd-layer package (cli/system/cmd/statusline) owns mode
// decisions (enabled, show_cost) and output; this core package owns
// the pure rendering: payload modeling, location resolution, path
// abbreviation, and printable-ASCII sanitization.
//
// All rendering degrades instead of failing: missing payload fields
// yield empty segments, and hostile content is stripped to bounded
// printable ASCII so it cannot corrupt the terminal
// (specs/statusline.md).
package statusline
