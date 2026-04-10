//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package fix

// Result tracks fixes applied during drift fix.
//
// Fields:
//   - Fixed: Number of issues successfully fixed
//   - Skipped: Number of issues skipped (not auto-fixable)
//   - Errors: Error messages from failed fix attempts
type Result struct {
	Fixed   int
	Skipped int
	Errors  []string
}
