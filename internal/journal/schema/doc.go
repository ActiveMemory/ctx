//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package schema validates Claude Code JSONL session files.
//
// Claude Code stores sessions as JSONL files with an undocumented,
// unversioned format that changes across releases. This package
// defines the expected record shape (known fields, record types,
// content block types) derived from empirical analysis, and
// validates raw lines against it to detect drift.
//
// Validation is strictly informational: it accumulates findings
// into a Collector but never blocks imports or other operations.
// Findings include unknown fields, missing required fields,
// unknown record types, and unknown content block types.
package schema
