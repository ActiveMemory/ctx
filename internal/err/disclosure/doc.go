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
// They fall into three groups:
//
//   - **Structure** ([ErrMultipleThemes], [ErrEntryBelowThemes],
//     [ErrStagingUnparsable]): the precondition refused a malformed root.
//   - **Cross-file** ([ErrOrphanThemeFile], [ErrMissingThemeFile],
//     [ErrDuplicateEntry], [ErrBrokenThemeLink]): the root ↔ theme-file
//     link graph or the one-place-per-entry invariant is broken.
//   - **Mover** ([ErrApplyNotEntryKind], [ErrEmptyAssignment],
//     [ErrEntryAssignedTwice], [ErrEntryNotInStaging], [ErrVerifyFailed]):
//     the milestone-3 digesting pass refused a malformed plan, was handed
//     an unsupported kind, or aborted because a moved body was not
//     byte-present after its theme-file append.
//
// # Wrapping strategy
//
// The structure and cross-file guards are returned bare — the guard
// identity is what the checks assert. The mover sentinels are likewise
// bare; the operator-facing message is resolved from errors.yaml at the
// CLI boundary.
//
// # Concurrency
//
// Package-level constants; safe for concurrent use.
package disclosure
