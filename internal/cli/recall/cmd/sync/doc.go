//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sync implements the ctx recall sync subcommand.
//
// It scans journal markdowns and syncs their frontmatter lock state
// into the state file, treating frontmatter as the source of truth.
package sync
