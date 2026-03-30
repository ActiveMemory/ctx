//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sync implements the ctx recall sync subcommand.
//
// [Cmd] builds the cobra.Command. [Run] scans journal Markdown
// files and updates .state.json to match each file's frontmatter
// lock status — the inverse of ctx recall lock.
package sync
