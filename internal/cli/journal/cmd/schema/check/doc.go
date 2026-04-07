//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package check provides the journal schema check subcommand.
//
// It walks JSONL session files in Claude Code project directories,
// validates each line against the embedded schema, and reports
// unknown fields, missing required fields, unknown record types,
// and unknown content block types. When drift is found, a
// Markdown report is written to .context/reports/schema-drift.md.
// When drift resolves, the report is automatically deleted.
// Exit code 0 means clean; exit code 1 means drift detected.
package check
