//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package evidence supplies the sentinel-error message
// strings, wrapping-format strings, and rendering tokens used
// by [github.com/ActiveMemory/ctx/internal/err/kb/evidence] and
// [github.com/ActiveMemory/ctx/internal/write/kb/evidence].
//
// Hosting these literals here honours the project-wide rule
// that magic strings belong under `internal/config/`; the err
// package consumes the error constants while the writer
// consumes the markdown rendering constants.
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/err/kb/evidence]
//     declares the sentinel error vars + constructors that
//     consume these strings.
//   - [github.com/ActiveMemory/ctx/internal/write/kb/evidence]
//     is the primary caller-side surface.
package evidence
