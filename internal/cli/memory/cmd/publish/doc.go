//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package publish implements the "ctx memory publish"
// command.
//
// # Overview
//
// The publish command selects high-value context from
// .context/ files, formats it as a marked block, and
// writes it into MEMORY.md. This makes curated project
// context available to AI tools that read MEMORY.md
// but do not have direct access to .context/.
//
// The published block is delimited by markers so it can
// be updated or removed by subsequent publish or
// unpublish operations without affecting user-authored
// content in MEMORY.md.
//
// # Flags
//
//	--budget <n>   Maximum line count for the published
//	               block (default from config).
//	--dry-run      Show what would be published without
//	               writing any files.
//
// # Behavior
//
// [Cmd] builds the cobra.Command and registers the two
// flags. [Run] discovers MEMORY.md, selects content
// from the context directory up to the budget limit,
// prints a plan showing counts by category (tasks,
// decisions, conventions, learnings) and total lines,
// then either writes the block (normal mode) or stops
// after the plan (dry-run mode).
//
// If MEMORY.md cannot be discovered, the command prints
// a warning and returns a "not found" error.
//
// # Output
//
// Prints a publication plan with category counts and
// total lines, followed by a "done" confirmation or
// a "dry run" notice.
package publish
