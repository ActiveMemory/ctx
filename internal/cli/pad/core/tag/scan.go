//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package tag

import "github.com/ActiveMemory/ctx/internal/cli/pad/core/blob"

// ScanText returns the portion of an entry to scan for tags.
// For blob entries, only the label is returned. For plain entries,
// the full text is returned.
//
// Parameters:
//   - entry: Scratchpad entry string
//
// Returns:
//   - string: Text to scan for tags
func ScanText(entry string) string {
	if label, _, ok := blob.Split(entry); ok {
		return label
	}
	return entry
}
