//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sourcecoverage supplies the sentinel-error
// messages, wrapping-format strings, and rendering tokens
// used by
// [github.com/ActiveMemory/ctx/internal/err/kb/sourcecoverage]
// and
// [github.com/ActiveMemory/ctx/internal/write/kb/sourcecoverage].
//
// Hosting these literals here honours the project-wide rule
// that magic strings belong under `internal/config/`.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/err/kb/sourcecoverage]
//     declares the sentinel error vars + constructors.
//   - [github.com/ActiveMemory/ctx/internal/write/kb/sourcecoverage]
//     is the primary caller-side surface.
package sourcecoverage
