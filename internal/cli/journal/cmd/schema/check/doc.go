//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package check implements the "ctx journal schema check"
// command.
//
// # Overview
//
// The check command scans JSONL session files in Claude
// Code project directories and validates each line against
// the embedded schema definition. It detects unknown
// fields, missing required fields, unknown record types,
// and unrecognized content block types.
//
// # Flags
//
//	--dir            Scan a specific directory instead of
//	                 the default Claude project paths.
//	--all-projects   Scan all discovered project
//	                 directories, not just the current one.
//	-q, --quiet      Suppress normal output; exit code
//	                 alone indicates pass (0) or drift (1).
//
// # Output
//
// When no drift is found, prints a clean summary with
// the number of files and lines scanned. When drift is
// detected, prints a categorized drift summary to stdout
// and writes a Markdown report to
// .context/reports/schema-drift.md. When drift later
// resolves, the report file is automatically deleted.
//
// In quiet mode, the command produces no output and
// relies solely on the exit code.
//
// # Behavior
//
// [Cmd] builds the cobra.Command and registers the three
// flags above. [Run] calls the core schema checker,
// optionally writes the report, and formats output based
// on drift status and the quiet flag.
//
// Designed for use in CI pipelines, nightly cron jobs,
// and interactive troubleshooting.
package check
