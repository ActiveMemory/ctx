//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package statusline holds display constants for the
// "ctx system statusline" renderer (specs/statusline.md).
//
// The status line is informational only: it renders model, context
// usage, and session cost from the payload Claude Code provides.
// There is deliberately no cost gating and no model-switch nudging;
// the spec's Decisions section records why.
//
// Constants are grouped as rendering limits (segment/line/path
// bounds), segment formats and separators, path display tokens, and
// the printable-ASCII sanitization range.
package statusline
