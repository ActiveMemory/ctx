//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

// Reminder represents a single session-scoped reminder.
//
// Fields:
//   - ID: Auto-incremented reminder identifier
//   - Message: Reminder text
//   - Created: ISO 8601 creation timestamp
//   - After: Optional trigger date (YYYY-MM-DD), nil for immediate
type Reminder struct {
	ID      int     `json:"id"`
	Message string  `json:"message"`
	Created string  `json:"created"`
	After   *string `json:"after"` // nullable YYYY-MM-DD
}
