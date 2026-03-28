//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package normalize

// TurnMatch holds the result of matching a turn header line.
type TurnMatch struct {
	Num  int
	Role string
	Time string
}
