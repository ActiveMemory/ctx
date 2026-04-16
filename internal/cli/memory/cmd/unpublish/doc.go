//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package unpublish implements the "ctx memory unpublish"
// command.
//
// # Overview
//
// The unpublish command removes the ctx-managed marker
// block from MEMORY.md that was previously written by
// "ctx memory publish". All user-authored content
// outside the markers is preserved intact.
//
// This is the inverse of the publish command and is
// useful when the user wants to stop sharing curated
// context through MEMORY.md or wants to rewrite the
// published section manually.
//
// # Flags
//
// This command accepts no flags.
//
// # Behavior
//
// [Cmd] builds a simple cobra.Command with no flags.
// [Run] discovers the source MEMORY.md, reads its
// content, and calls memory.RemovePublished to strip
// the marked block. If no published block is found,
// prints a "not found" message and returns nil. If
// the block is found, writes the cleaned content back
// to the file and prints a confirmation.
//
// If the source MEMORY.md cannot be discovered, the
// command prints a warning and returns a "not found"
// error.
//
// # Output
//
// Prints either a "not found" notice when no published
// block exists, or an "unpublished" confirmation when
// the block was successfully removed.
package unpublish
