//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package persistence

// State holds the counter state for persistence nudging.
//
// Fields:
//   - Count: Edit/Write calls since last nudge
//   - LastNudge: Prompt number of the last nudge
//   - LastMtime: Unix timestamp of last TASKS.md modification
type State struct {
	Count     int
	LastNudge int
	LastMtime int64
}
