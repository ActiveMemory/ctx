//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package dream defines the typed error constructors returned by
// [internal/dream]: guard refusals (write-scope, don't-leak),
// state-file persistence failures, ledger append/read failures, and
// proposal-validation failures.
//
// # Why Typed Errors
//
//   - Stability: error categories are part of the public API.
//   - Routing: messages are sourced from the YAML text registry via
//     [internal/assets/read/desc], keyed by DescKey constants in
//     [internal/config/embed/text].
//   - Wrapping: constructors that take a cause wrap it via %w so
//     callers can errors.Is/errors.As against the underlying error.
//
// # The Guard Refusals
//
// [WriteScope] and [Leak] are the structural safety invariants of the
// dream rendered as errors: a write target outside dreams/ or ideas/
// (and specs/ only via an accepted promote) is refused, and a target
// that resolves to a git-tracked path is refused. They are the
// load-bearing portability requirement of the executor contract.
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package dream
