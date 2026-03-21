//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package core

// Reminder represents a single session-scoped reminder.
type Reminder struct {
	ID      int     `json:"id"`
	Message string  `json:"message"`
	Created string  `json:"created"`
	After   *string `json:"after"` // nullable YYYY-MM-DD
}
