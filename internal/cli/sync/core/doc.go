//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core holds the **shared helpers** that the
// `ctx sync` subcommands all rely on: the action
// resolver that turns a scanner finding into a typed
// "consider documenting this" suggestion, the renderer
// that formats the resulting list, and the predicate
// helpers that decide what counts as "undocumented".
//
// # Public Surface
//
// Each suggestion produced by the scanner is an
// [Action] with a kind (package-file / config /
// directory), a path, and a one-line "consider
// documenting in X" pointer. The helpers here build
// the slice; the CLI surface ([internal/cli/sync])
// orchestrates the run and the rendering.
//
// # Sub-Packages
//
//   - **[validate]** — the predicate package: the
//     type-aware "is this an undocumented X"
//     checks.
//
// # Concurrency
//
// Filesystem-bound and stateless. Concurrent
// invocations against the same project each pay
// the full scan cost.
//
// # Related Packages
//
//   - [internal/cli/sync]               — the
//     `ctx sync` CLI surface.
//   - [internal/cli/sync/core/validate] — the
//     undocumented-artifact predicates.
//   - [internal/assets/read/lookup]     — supplies
//     the config-file pattern set the predicates
//     consult.
package core
