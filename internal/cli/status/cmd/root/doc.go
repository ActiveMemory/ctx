//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the **`ctx status`** command,
// which prints a summary of the current project context.
//
// # What It Does
//
// The command loads the resolved .context/ directory and
// displays the state of every context file: TASKS.md,
// DECISIONS.md, LEARNINGS.md, CONVENTIONS.md, and any
// other tracked files. It reports file existence, size,
// and last-modified timestamps so the user can quickly
// gauge how fresh and complete their context is.
//
// # Flags
//
//   - **--json**: emit the status report as a single
//     JSON object for machine consumption (scripts,
//     dashboards, CI).
//   - **--verbose, -v**: include inline previews of
//     each file's opening lines so the user can skim
//     content without opening individual files.
//
// # Output
//
// In default (text) mode the output is a human-readable
// table with one row per context file. In JSON mode
// the output is a single JSON object written to stdout.
// Both modes respect the --verbose flag to optionally
// include file content previews.
//
// # Delegation
//
// [Cmd] builds the cobra command and binds flags.
// [Run] loads the context via [context/load.Do] and
// delegates rendering to [cli/status/core/out], which
// owns the text and JSON formatters. If the .context/
// directory does not exist, Run returns a user-facing
// "context not initialized" error.
package root
