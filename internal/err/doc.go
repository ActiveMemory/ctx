//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package err is the root of the typed-error tree for
// the ctx CLI. Each child package under err/ defines
// domain-scoped error constructors for a single
// subsystem (e.g. [err/backup], [err/config],
// [err/journal]).
//
// # Design
//
// The err/ tree separates error construction from
// error handling. Constructors live here; renderers
// live in internal/write/*. This split gives three
// benefits:
//
//   - **Stable categories** -- adding a new error
//     constructor is an explicit, reviewable change.
//   - **Localized text** -- all user-facing wording
//     is looked up through [internal/assets/read/desc]
//     so that error messages are consistent and
//     centrally maintained.
//   - **Wrapping** -- every constructor that accepts a
//     cause wraps it with fmt.Errorf %w, so callers
//     can use errors.Is / errors.As against system
//     errors (os.ErrNotExist, io.EOF, etc.).
//
// # Package Layout
//
// Each child package exports:
//
//   - Pure constructor functions (no state, no IO).
//   - Occasionally a sentinel error variable (e.g.
//     [err/schema.ErrDrift]).
//   - Occasionally a typed error struct (e.g.
//     [err/context.NotFoundError]).
//
// # Concurrency
//
// All constructors are pure functions. Concurrent
// callers never race.
package err
