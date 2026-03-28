//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package persistence

// State holds the counter state for persistence nudging.
type State struct {
	Count     int
	LastNudge int
	LastMtime int64
}
