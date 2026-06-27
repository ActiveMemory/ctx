//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hubsync

import "time"

// SetSyncTimeoutForTest overrides the session-start pull deadline
// and returns a restore func. It lets external tests drive the
// deadline against an unresponsive hub in milliseconds instead of
// the full production interval.
func SetSyncTimeoutForTest(d time.Duration) func() {
	prev := syncTimeout
	syncTimeout = d
	return func() { syncTimeout = prev }
}
