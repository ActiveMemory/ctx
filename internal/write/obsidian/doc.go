//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package obsidian provides terminal output for the
// Obsidian vault generation command (ctx journal
// obsidian).
//
// # Exported Functions
//
// [InfoGenerated] reports the number of journal entries
// written and the output directory path after vault
// generation completes. It also prints a "Next Steps"
// section with instructions for opening the vault in
// Obsidian, including the exact directory path to use.
//
// # Message Categories
//
//   - Info: generation result with entry count and
//     output path, followed by next-step guidance
//
// # Usage
//
//	obsidian.InfoGenerated(cmd, entryCount, outputDir)
//
// The output includes both the result confirmation
// and actionable next steps so the user can immediately
// open the generated vault without consulting docs.
package obsidian
