//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package validate

import (
	errPad "github.com/ActiveMemory/ctx/internal/err/pad"
)

// Index checks that n is a valid 1-based index into entries.
//
// Parameters:
//   - n: 1-based entry index
//   - entries: The entries to validate against
//
// Returns:
//   - error: Non-nil if n is out of range
func Index(n int, entries []string) error {
	if n < 1 || n > len(entries) {
		return errPad.EntryRange(n, len(entries))
	}
	return nil
}
