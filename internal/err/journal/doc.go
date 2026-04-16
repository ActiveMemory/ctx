//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package journal defines the **typed error constructors**
// returned by the journal pipeline — state-file load
// and save failures, missing journal directories, and
// related operational errors that reach the user
// through `ctx journal *` subcommands.
//
// # Why Typed Errors
//
//   - **Stability** — error categories are part of
//     the public API.
//   - **Routing** — write-side packages map error
//     types to localized text via
//     [internal/assets/read/desc].
//   - **Wrapping** — constructors wrap the
//     underlying cause via `%w` so callers can
//     `errors.Is` against system errors when
//     needed.
//
// # Public Surface
//
// Constructors: [LoadState], [SaveState],
// [LoadStateErr], [LoadStateFailed],
// [SaveStateFailed], [NoDir].
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package journal
