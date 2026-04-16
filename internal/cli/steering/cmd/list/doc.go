//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package list implements the **`ctx steering list`**
// subcommand, which displays all steering files in the
// project steering directory.
//
// # What It Does
//
// The command loads every .md file from the steering
// directory, parses its YAML frontmatter, and prints a
// one-line summary for each file containing:
//
//   - name
//   - inclusion mode (always, auto, manual)
//   - priority (integer sort key)
//   - target tools (or "all tools" when unscoped)
//
// After the file entries, a total count line is printed.
// If no steering files exist, a "no files found"
// message is displayed instead.
//
// # Arguments
//
// None. The command accepts no positional arguments
// (cobra.NoArgs).
//
// # Output
//
// A line-per-file table followed by a count footer.
// Example:
//
//	product   always  10  all tools
//	ci-hints  manual  50  claude, copilot
//	3 steering files
//
// # Delegation
//
// [Cmd] builds the cobra command. [Run] calls
// [steering.LoadAll] to parse the directory, then
// iterates over results and emits each entry through
// [write/steering.FileEntry]. The final count is
// emitted via [write/steering.FileCount].
package list
