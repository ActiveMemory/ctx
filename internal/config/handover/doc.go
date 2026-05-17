//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package handover supplies the sentinel-error message strings,
// wrapping-format strings, and section-header constants used by
// the handover writer (Phase KB).
//
// Holding the strings here honours the project-wide rule that
// magic strings belong under `internal/config/`; the
// [github.com/ActiveMemory/ctx/internal/err/handover] package
// consumes these constants when constructing user-facing
// errors, and the
// [github.com/ActiveMemory/ctx/internal/write/handover]
// package consumes the section-header constants when composing
// the markdown body.
package handover
