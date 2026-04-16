//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package trace defines the typed error constructors
// for the `ctx trace` command, which manages git
// commit annotations, hook installation, and history
// recording for context-aware commit messages.
//
// # Domain
//
// Errors fall into four categories:
//
//   - **Git operations**: git rev-parse or git
//     log failed, or a commit ref could not be
//     resolved. Constructors: [GitDir], [GitLog],
//     [ResolveCommit].
//   - **Hook management**: a non-ctx hook already
//     exists, or writing the hook script failed.
//     Constructors: [HookExists], [HookWrite].
//   - **History / override IO**: writing the
//     trace history or override file failed.
//     Constructors: [WriteHistory],
//     [WriteOverride].
//   - **Validation**: the --note flag is missing,
//     or an unknown action was provided.
//     Constructors: [NoteRequired],
//     [UnknownAction].
//
// # Wrapping Strategy
//
// IO constructors wrap their cause with
// fmt.Errorf %w so callers can inspect the
// underlying error. [NoteRequired] returns a
// plain errors.New value. All user-facing text
// is resolved through [internal/assets/read/desc].
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package trace
