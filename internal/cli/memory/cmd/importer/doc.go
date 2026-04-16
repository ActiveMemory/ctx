//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package importer implements the "ctx memory import"
// command.
//
// # Overview
//
// The import command parses entries from MEMORY.md,
// classifies each one by heuristic keyword matching,
// deduplicates against previously imported entries, and
// promotes new entries into the appropriate .context/
// files (TASKS.md, DECISIONS.md, CONVENTIONS.md, or
// LEARNINGS.md).
//
// # Flags
//
//	--dry-run    Show the classification plan without
//	             writing any files. Each entry is
//	             printed with its target file and the
//	             keywords that triggered classification.
//
// # Behavior
//
// [Cmd] builds the cobra.Command and registers the
// --dry-run flag. [Run] performs the following steps:
//
//  1. Discovers the source MEMORY.md in the project.
//  2. Parses all entries from the file content.
//  3. Loads import state to identify duplicates.
//  4. For each entry, computes a content hash, skips
//     already-imported entries, classifies the entry
//     by keyword matching, and either promotes it
//     (normal mode) or reports it (dry-run mode).
//  5. Saves updated import state with hashes of newly
//     imported entries.
//
// Entries classified as "skip" are silently ignored
// unless --dry-run is active.
//
// # Output
//
// Prints a scan header with the source name and entry
// count, followed by per-entry results (added or
// classified), and a summary with counts by target file
// plus duplicates and skipped entries.
package importer
