//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package index renders the output of the `ctx index` command: the headings
// projected from a knowledge file, either as indented lines (the default,
// deeper levels indented under their parents) or as a JSON array of
// {level, text} objects (--json).
//
// It is the write-side counterpart to the projection logic in
// [internal/heading]: the recognizer extracts headings, this package emits
// them. Keeping emission here honors the convention that CLI output routes
// through internal/write, not through cmd/ or core/ directly.
package index
