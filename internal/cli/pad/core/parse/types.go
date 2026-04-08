//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parse

// Entry represents a pad entry with a stable ID.
//
// Fields:
//   - ID: Stable auto-incrementing identifier
//   - Content: Entry content without the ID prefix
type Entry struct {
	ID      int
	Content string
}
