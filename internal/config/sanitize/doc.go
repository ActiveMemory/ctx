//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sanitize defines string and length constants used by
// the sanitize layer.
//
// Constants are referenced by internal/sanitize via config/sanitize.*.
// Provides: [NullByte], [DotDot], [ForwardSlash], [Backslash],
// [HyphenReplace], [EscapePrefix], [MaxSessionIDLen].
package sanitize
