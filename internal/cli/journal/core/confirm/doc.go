//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package confirm handles user confirmation prompts for
// journal import operations.
//
// # Overview
//
// Before the journal import command writes files to
// disk, it presents a summary of what will happen and
// asks the user to confirm. This package implements
// that interactive confirmation step.
//
// # Behavior
//
// [Import] renders the import plan summary (new, regenerated,
// skipped, locked file counts), prompts the user for
// confirmation on stdin, and returns true only when the
// response is "y" or "yes".
//
// # Data Flow
//
// When [Import] is called it:
//
//  1. Renders the import plan summary showing counts
//     for new files, regenerated files, skipped files,
//     and locked files via writeRecall.ImportSummary.
//  2. Prints a confirmation prompt via
//     writeRecall.ConfirmPrompt.
//  3. Reads a line from stdin using a buffered reader.
//  4. Trims whitespace and lowercases the response.
//  5. Returns true if the response matches "y" or
//     "yes", false otherwise.
//  6. Returns an error if reading from stdin fails.
package confirm
