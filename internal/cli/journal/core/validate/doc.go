//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package validate provides input validation for
// journal commands.
//
// Before the journal import pipeline begins planning
// or writing, inputs must be checked for consistency.
// This package contains the validation functions that
// the cmd/journal layer calls early in the import
// flow.
//
// # Message Validation
//
// [EmptyMessage] checks whether a session message has
// no substantive content. A message is empty when its
// Text field is blank and it carries no ToolUses or
// ToolResults. The plan package uses EmptyMessage to
// filter out empty messages before splitting sessions
// into parts, ensuring that journal files contain only
// meaningful conversation turns.
//
// # Flag Validation
//
// [ImportFlags] validates the combination of import
// flags and positional arguments. Two rules are
// enforced:
//
//   - Passing a session ID together with --all is
//     an error (AllWithID).
//   - Using --regenerate without --all is an error
//     (RegenerateRequiresAll).
//
// Both checks return typed errors from the err/session
// and err/journal packages, which the cmd layer
// renders to the user.
//
// # Data Flow
//
// The cmd/journal layer calls ImportFlags before
// calling query.FindSessions. EmptyMessage is called
// by plan.Import during action planning. No state is
// mutated; both functions are pure predicates.
package validate
