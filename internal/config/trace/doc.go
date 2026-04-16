//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package trace defines constants for the ctx trace
// subsystem, which embeds context references into git
// commit messages and tracks commit history.
//
// ctx trace installs git hooks that automatically
// append a ctx-context trailer to every commit,
// linking commits to decisions, learnings, sessions,
// and tasks. This package provides the trailer
// format, reference type identifiers, display
// defaults, diff parsing tokens, and storage
// filenames.
//
// # Git Trailer
//
//   - [TrailerKey] ("ctx-context") — the git trailer
//     key embedded in commit messages.
//   - [TrailerFormat] — the format string for the
//     trailer line.
//
// # Reference Types
//
//   - [RefTypeNote], [RefTypeSession],
//     [RefTypeDecision], [RefTypeLearning],
//     [RefTypeConvention], [RefTypeTask] — identifiers
//     used in ctx-context trailer values.
//   - [RefFormat], [SessionRefFormat] — format strings
//     for numbered and session refs.
//
// # Display Defaults
//
//   - [DefaultLastFile] (20) — commits shown by
//     ctx trace file.
//   - [DefaultLastShow] (10) — commits shown by
//     ctx trace with no arguments.
//   - [ShortHashLen] (7) — abbreviated hash length.
//
// # Hook Management
//
//   - [ActionEnable], [ActionDisable] — arguments
//     for ctx trace hook enable/disable.
//   - [CtxTraceMarker] — sentinel string that
//     identifies ctx-installed git hooks.
//   - [ScriptPrepareCommitMsg],
//     [ScriptPostCommit] — embedded hook script
//     filenames.
//
// # Storage Files
//
//   - [FileHistory], [FileOverrides],
//     [FilePending] — JSONL filenames within the
//     trace state directory.
//
// # Concurrency
//
// All exports are immutable. Safe for any access
// pattern.
package trace
