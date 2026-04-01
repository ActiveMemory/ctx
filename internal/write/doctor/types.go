//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package doctor

// ResultItem holds the display data for a single doctor check result.
//
// Fields:
//   - Category: Check category (e.g. "context", "hooks")
//   - Status: Pass/fail indicator symbol
//   - Message: Human-readable result description
type ResultItem struct {
	Category string
	Status   string
	Message  string
}
