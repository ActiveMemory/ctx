//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package data

// DocEntry pairs a document alias with its display label.
//
// Fields:
//   - Alias: Document lookup key (e.g. "manifesto")
//   - Label: Menu display text
type DocEntry struct {
	Alias string
	Label string
}
