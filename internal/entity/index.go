//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entity

// IndexEntry represents a parsed entry header from a context
// file.
//
// Fields:
//   - Timestamp: Full timestamp (YYYY-MM-DD-HHMMSS)
//   - Date: Date only (YYYY-MM-DD)
//   - Title: Entry title
type IndexEntry struct {
	Timestamp string
	Date      string
	Title     string
}
