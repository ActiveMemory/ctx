//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package obsidian builds an Obsidian vault from journal entries.
//
// [BuildVault] handles the full file generation pipeline:
// scan entries, create directories, transform frontmatter,
// convert links, build MOC pages, and write Home.md.
package obsidian
