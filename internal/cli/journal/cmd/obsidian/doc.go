//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package obsidian implements the ctx journal obsidian subcommand.
//
// [Cmd] builds the cobra.Command with --output flag. [Run]
// delegates to core/obsidian.BuildVault to generate an
// Obsidian vault from journal entries.
package obsidian
