//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package disclosure holds the identity sentinels for the
// progressive-disclosure guards and invariants (see
// specs/progressive-disclosure.md).
//
// # Domain
//
// Each sentinel is an [github.com/ActiveMemory/ctx/internal/entity.Sentinel]
// — a string const whose user-facing text is resolved from
// commands/text/errors.yaml at call time. Callers match them with
// errors.Is against the value returned by [Validate] and the invariant
// checks in internal/disclosure.
//
// They fall into two groups:
//
//   - **Structure** ([ErrMultipleThemes], [ErrEntryBelowThemes],
//     [ErrStagingUnparsable]): the precondition refused a malformed root.
//   - **Cross-file** ([ErrOrphanThemeFile], [ErrMissingThemeFile],
//     [ErrDuplicateEntry], [ErrBrokenThemeLink]): the root ↔ theme-file
//     link graph or the one-place-per-entry invariant is broken.
//
// # Wrapping strategy
//
// Milestone 1 returns these bare — the guard identity is what the checks
// assert. Parameterized wrappers that name the specific theme file or
// entry are deferred to the digesting pass (a later milestone), which is
// where the message reaches an operator.
//
// # Concurrency
//
// Package-level constants; safe for concurrent use.
package disclosure
