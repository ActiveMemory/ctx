//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package memory defines the **typed error constructors**
// returned by [internal/memory] — the Claude Code auto-
// memory bridge — for discovery, diff, mirror, and
// publish failures.
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
// Constructors: [NotFound], [DiscoverFailed],
// [DiffFailed], [SelectContentFailed],
// [PublishFailed], [Read].
//
// # Why "NotFound" Is Distinct from "DiscoverFailed"
//
// "Auto memory does not exist for this project"
// ([NotFound]) is a normal state Claude Code
// returns for projects with no recorded memory;
// the CLI surfaces it as "no memory yet, run a
// session first". Discover failures (path
// resolution errors, permission denied) are real
// errors that need user attention.
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package memory
