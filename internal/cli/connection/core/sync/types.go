//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

// state tracks the last synced sequence number.
//
// Fields:
//   - LastSequence: hub sequence of the most recent entry
type state struct {
	LastSequence uint64 `json:"last_sequence"`
}
