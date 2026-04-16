//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package decision implements the ctx decision command
// group for managing DECISIONS.md.
//
// DECISIONS.md records architectural choices with their
// rationale, trade-offs, and timestamps. The decision
// command group provides tooling to maintain this file's
// quick-reference index table, which maps decision
// numbers to one-line summaries for fast scanning.
//
// # Subcommands
//
//   - reindex: scans DECISIONS.md entries and regenerates
//     the index table at the top of the file, ensuring
//     numbering and summaries stay consistent with the
//     full entries below
//
// # Subpackages
//
//	cmd/reindex -- cobra command for index regeneration
package decision
