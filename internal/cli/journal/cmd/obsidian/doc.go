//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package obsidian implements the ctx journal obsidian subcommand.
//
// [Cmd] builds the cobra.Command with --output flag. [Run]
// generates an Obsidian vault from journal entries with wikilinks,
// topic pages, and frontmatter. [BuildVault] handles the file
// generation pipeline.
package obsidian
