//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package reminder defines the typed error
// constructors for session reminder operations.
// These errors fire when reading, parsing, or
// looking up reminders stored in the context
// directory.
//
// # Domain
//
// Errors fall into two categories:
//
//   - **File IO** -- the reminders file could not
//     be read or parsed. Constructors: [Read],
//     [Parse].
//   - **Lookup** -- no reminder matches the given
//     ID, or no ID was provided. Constructors:
//     [NotFound], [IDRequired].
//
// # Wrapping Strategy
//
// [Read] and [Parse] wrap their cause with
// fmt.Errorf %w so callers can inspect the
// underlying error. [NotFound] returns a plain
// formatted error. [IDRequired] returns a plain
// errors.New value. All user-facing text is
// resolved through [internal/assets/read/desc].
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package reminder
