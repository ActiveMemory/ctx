//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package wrap soft-wraps long lines in markdown files to a target
// width (default 80 characters).
//
// [Content] wraps all lines in a journal entry. [ContextFile] wraps
// lines in a context file (.context/*.md), handling markdown list
// continuation with 2-space indent. [Soft] wraps a single line at
// word boundaries, returning multiple lines.
package wrap
