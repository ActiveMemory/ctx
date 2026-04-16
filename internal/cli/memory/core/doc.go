//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core provides shared helpers for the memory
// subcommands.
//
// The "ctx memory" command family reports statistics
// about context files such as line counts and token
// estimates. This core package holds the counting
// logic that the cmd/memory layer delegates to.
//
// # Line Counting
//
// The count sub-package exports [count.FileLines],
// which counts the number of newline characters in a
// byte slice. It uses bytes.Count with the LF token
// from config/token. The cmd/memory layer reads each
// context file into memory and passes the raw bytes
// to FileLines to obtain per-file line counts for
// display.
//
// # Data Flow
//
// The cmd/memory layer discovers context files, reads
// them from disk, and calls FileLines for each file.
// The resulting counts are passed to the write/memory
// package for formatted output showing file sizes and
// totals.
package core
