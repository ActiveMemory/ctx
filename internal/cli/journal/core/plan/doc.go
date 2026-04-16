//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package plan builds the import plan for journal import
// operations.
//
// When the user runs "ctx journal source", sessions must
// be matched against existing journal files to decide
// what to create, regenerate, skip, or rename. This
// package performs that planning without writing any
// files.
//
// # Import Function
//
// [Import] is the sole exported function. It accepts a
// list of sessions, the journal output directory, an
// index mapping session IDs to existing filenames, the
// journal processing state, import flags, and a flag
// indicating single-session mode.
//
// # Planning Algorithm
//
// For each session, Import:
//
//  1. Filters out empty messages (no text, tool uses,
//     or tool results).
//  2. Splits long sessions into multiple parts using
//     MaxMessagesPerPart from the journal config.
//  3. Resolves a title-based slug by checking the
//     existing file for a frontmatter title, falling
//     back to slug generation.
//  4. Detects renames when an old slug differs from
//     the new one, recording a RenameOp.
//  5. Assigns each part an action: ActionNew for files
//     that do not exist, ActionLocked for entries
//     locked in state or frontmatter, ActionRegenerate
//     when forced by flags, or ActionSkip otherwise.
//
// The result is an ImportPlan containing all FileAction
// items, RenameOps, and counters (NewCount, LockedCount,
// RegenCount, SkipCount).
//
// # Connection to Other Layers
//
// The cmd/journal package calls [Import] and then hands
// the plan to the write/journal package for execution.
// Locked entries are never overwritten; frontmatter
// locks are promoted to persistent state so future
// operations skip re-parsing.
package plan
