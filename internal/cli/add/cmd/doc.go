//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package cmd wires the cobra subcommands for ctx add.
//
// It registers decision, learning, convention, and task subcommands
// under the add parent, following the cmd/root + core taxonomy.
// Each subcommand delegates to the shared [root.Run] function
// with type-specific flags.
package cmd
