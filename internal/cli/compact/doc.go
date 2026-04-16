//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package compact implements the "ctx compact" command for
// cleaning up and consolidating context files.
//
// The compact command performs maintenance on .context/
// files by moving completed tasks to a dedicated archive
// section, removing empty sections, and optionally
// archiving old content that exceeds configured retention
// thresholds. This keeps the active context lean so that
// AI agents can consume it within token budgets.
//
// # What Compact Does
//
// When invoked, compact walks TASKS.md and relocates
// checked-off items to a "Completed" section at the
// bottom of the file. Empty headings left behind are
// pruned. Other context files (DECISIONS.md, LEARNINGS.md)
// can be compacted to remove superseded entries.
//
// # Subpackages
//
//   - cmd/root: cobra command definition and flag binding
//   - core: task relocation, section pruning, and
//     archival logic
package compact
