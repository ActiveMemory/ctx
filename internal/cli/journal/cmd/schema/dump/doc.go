//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package dump provides the journal schema dump subcommand.
//
// It prints the embedded JSONL schema definition to stdout,
// showing all known record types with their required and optional
// fields, and all recognized content block types with their
// parse status. The output is human-readable and useful for
// understanding what the schema validator expects before
// running a check.
package dump
