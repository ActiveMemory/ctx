//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package status implements the "ctx status" command for
// displaying context health and summary information.
//
// The status command reads context files from .context/
// and displays a point-in-time snapshot of the project's
// context state. This is typically the first command an
// AI agent or user runs to orient themselves.
//
// # Output Includes
//
//   - File presence: which context files exist
//   - Token estimates: approximate token cost of each
//     file, useful for budget planning
//   - Modification times: when each file was last changed
//   - Task completion ratio: checked vs total tasks in
//     TASKS.md
//   - Entry counts: number of decisions, learnings,
//     conventions
//
// Output can be in human-readable (default) or JSON
// format via --json.
//
// # Subpackages
//
//   - cmd/root: cobra command definition and flag binding
//   - core: file scanning, token estimation, and
//     summary computation
package status
