//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package wrap soft-wraps long lines in journal files to
// approximately 80 characters.
//
// [Content] wraps all lines in a journal entry. [Soft] wraps a
// single line at word boundaries, returning multiple lines. Lines
// inside code fences are preserved as-is.
package wrap
