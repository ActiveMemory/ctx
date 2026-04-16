//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the **`ctx sync`** command,
// which reconciles the project context with the current
// state of the codebase.
//
// # What It Does
//
// The command scans the codebase for structural changes
// that should be reflected in context files: new
// directories, added or removed package manager files,
// configuration file changes, and similar discrepancies.
// Each detected issue is presented as a numbered action
// with a type label, description, and suggestion.
//
// # Flags
//
//   - **--dry-run** -- show what would change without
//     modifying any files. The output lists suggested
//     actions but does not prompt for confirmation.
//
// # Output
//
// When everything is in sync, prints an "all clear"
// message. Otherwise, prints a numbered list of
// detected actions followed by a summary count. In
// dry-run mode the summary reminds the user that no
// changes were applied.
//
// # Delegation
//
// [Cmd] builds the cobra command and binds the
// --dry-run flag. [Run] loads the context via
// [context/load.Do], calls [cli/sync/core/action.Detect]
// to find discrepancies, then emits results through the
// [write/sync] formatters. If the .context/ directory
// does not exist, Run returns a "context not
// initialized" error.
package root
