//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package statusline implements the hidden
// "ctx system statusline" cobra subcommand.
//
// Claude Code invokes the configured statusLine command after each
// assistant message and pipes a JSON payload to its stdin; the first
// line of stdout becomes the status line. ctx init wires this command
// into .claude/settings.local.json (see the initialize merge core).
//
// # Behavior
//
// The command renders one line assembled from up to four segments:
//
//	user@host dir | model | ctx: N% | $C.CC
//
// Location comes from the process environment plus
// workspace.current_dir (falling back to cwd), home-abbreviated.
// Model comes from model.display_name. Context percentage comes
// from context_window.used_percentage, which may be null early in
// a session. Cost comes from cost.total_cost_usd, rendered as a
// plain figure; the segment is suppressed when .ctxrc sets
// statusline.show_cost to false.
//
// Missing payload fields drop their segment rather than rendering a
// placeholder: a wrong-looking number is worse than an absent one.
// Malformed input degrades to whatever segments remain renderable.
// When statusline.enabled is false, the command prints an empty
// line so the display blanks immediately.
//
// The command always exits zero: Claude Code blanks the status line
// on a non-zero exit, so failures here must never be fatal. Output
// is sanitized to bounded printable ASCII; hostile payload content
// cannot corrupt the terminal.
//
// This surface is deliberately informational. There is no cost
// gating and no model-switch nudging (specs/statusline.md records
// the rationale).
//
// # Flags
//
// None. The command reads the status line JSON payload from stdin.
//
// # Output
//
// A single sanitized line on stdout.
//
// # Delegation
//
// Rendering primitives live in system/core/statusline; the line is
// written via write/statusline.
package statusline
