//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sync implements the "ctx journal sync" command.
//
// # Overview
//
// The sync command reconciles journal lock state by
// treating Markdown frontmatter as the source of truth.
// It scans all journal Markdown files and updates
// .state.json to match each file's frontmatter lock
// status. This is the inverse of "ctx journal lock",
// which writes state first and expects frontmatter to
// follow.
//
// # Flags
//
// This command accepts no flags.
//
// # Behavior
//
// [Cmd] builds a simple cobra.Command with no flags.
// [Run] performs the following steps:
//
//  1. Loads .state.json from the journal directory.
//  2. Discovers all journal Markdown files.
//  3. For each file, reads the frontmatter lock field
//     and compares it to the state file entry.
//  4. When frontmatter says locked but state says
//     unlocked, marks the entry as locked in state.
//  5. When frontmatter says unlocked but state says
//     locked, clears the lock in state.
//  6. Saves the updated .state.json.
//
// # Output
//
// Prints one line per state change (locked or unlocked)
// with the affected filename. Ends with a summary line
// showing total locked and unlocked counts. If no
// journal files are found, prints a "nothing to sync"
// message.
package sync
