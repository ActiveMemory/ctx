//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package bootstrap implements the
// **`ctx system bootstrap`** command, which emits
// context directory information for AI agents at
// session start.
//
// # What It Does
//
// The command resolves the .context/ directory path and
// prints it along with context file listings, rules,
// next steps, and plugin warnings. This is the first
// command an agent runs to anchor its session: it tells
// the agent where context lives and what to read next.
//
// This is agent-only plumbing; no human types it
// interactively.
//
// # Flags
//
//   - **--json**: emit a structured JSON object
//     containing the directory path, file list, rules,
//     next steps, and any plugin warnings.
//   - **--quiet, -q**: emit only the bare directory
//     path with no decoration, suitable for shell
//     variable capture.
//
// # Output
//
// In default (text) mode: a formatted block showing
// the context directory, a wrapped file list, numbered
// rules, numbered next steps, and an optional plugin
// warning. In --quiet mode: just the directory path.
// In --json mode: a single JSON object to stdout.
//
// # Delegation
//
// [Cmd] builds the cobra command and binds flags.
// [Run] stats the .context/ directory, collects file
// listings via [core/bootstrap.ListContextFiles],
// parses rules and next steps from embedded text
// assets, and delegates rendering to the
// [write/bootstrap] formatters.
package bootstrap
