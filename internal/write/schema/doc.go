//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package schema provides terminal output for the
// schema check and dump commands (ctx schema check,
// ctx schema dump).
//
// The schema system validates session state files
// against expected field layouts and detects drift.
// Output functions render validation results and
// raw schema dumps.
//
// # Validation Results
//
// [NoDirs] reports that no session directories were
// found. [NoFiles] reports that no session files
// were found within the directories. [Clean] prints
// a success message with the count of files and
// lines scanned when no drift is detected.
//
// [DriftSummary] prints a pre-formatted drift
// summary to stderr when validation finds
// mismatches between actual and expected schemas.
//
// # Schema Dump
//
// [DumpLine] prints a single line of schema dump
// output. [DumpBlank] prints a blank separator
// line between dump sections.
//
// Functions accept primitive types (strings, ints)
// rather than domain types to avoid cross-package
// type references.
package schema
