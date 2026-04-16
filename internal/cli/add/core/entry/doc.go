//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package entry provides type predicates and complexity
// detection for context file entries used by the add
// command.
//
// # Type Predicates
//
// Three predicate functions classify a user-supplied file
// type string into its canonical entry kind:
//
//   - [FileTypeIsTask] returns true when the input
//     resolves to a task entry (e.g. "task", "tasks").
//   - [FileTypeIsDecision] returns true for decision
//     entries (e.g. "decision", "decisions").
//   - [FileTypeIsLearning] returns true for learning
//     entries (e.g. "learning", "learnings").
//
// Each function delegates to config/entry.FromUserInput
// for alias resolution, so callers never deal with raw
// string matching.
//
// # Spec Nudge Detection
//
// [NeedsSpec] inspects task content to decide whether the
// add command should suggest creating a feature spec. It
// fires when the text exceeds the length threshold from
// .ctxrc (rc.SpecNudgeMinLen) or contains any of the
// design-signal words configured via rc.SpecSignalWords.
// The check is case-insensitive.
//
// # Data Flow
//
// The cmd/ layer passes the user-provided type string to
// these predicates to select the correct formatter in the
// format subpackage and the correct insertion strategy in
// the insert subpackage. NeedsSpec is called after content
// extraction to optionally emit a nudge message.
package entry
