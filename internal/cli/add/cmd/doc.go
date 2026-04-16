//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package cmd wires the cobra subcommands for the
// "ctx add" command tree.
//
// # Purpose
//
// This package registers the root add command and its
// type-specific subcommands under a single parent.
// Each subcommand delegates to [root.Run] with flags
// that control which context file receives the entry.
//
// # Subcommand Registration
//
// The package imports and wires the root subcommand
// from the root/ child package. The root command
// accepts a positional type argument (task, decision,
// learning, convention) and dispatches accordingly.
//
// # Entry Types
//
// Supported entry types and their target files:
//
//   - task       -> TASKS.md
//   - decision   -> DECISIONS.md
//   - learning   -> LEARNINGS.md
//   - convention -> CONVENTIONS.md
//
// Both singular and plural forms are accepted.
//
// # Output
//
// On success the command prints a confirmation line
// naming the file that was updated. When --share is
// set it also publishes the entry to a connected
// ctx Hub instance.
package cmd
