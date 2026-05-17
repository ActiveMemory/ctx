//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package closeout supplies the sentinel-error message strings
// and wrapping-format strings used by
// [github.com/ActiveMemory/ctx/internal/err/closeout]. Holding
// them here honours the project-wide rule that magic strings
// belong under `internal/config/`; the err package consumes
// these constants when constructing user-facing errors for the
// closeout writer (Phase KB).
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/err/closeout]
//     declares the sentinel error vars + constructors that
//     consume these strings.
//   - [github.com/ActiveMemory/ctx/internal/write/closeout]
//     is the primary caller-side surface; it imports the err
//     package, never these raw strings.
package closeout
