//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package normalize

// TurnMatch holds the result of matching a turn header line.
//
// Fields:
//   - Num: Turn number (1-based)
//   - Role: Speaker role (Human, Assistant)
//   - Time: Timestamp string from the header
type TurnMatch struct {
	Num  int
	Role string
	Time string
}
