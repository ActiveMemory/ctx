//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package architecture

// Map staleness hook configuration.
const (
	// MapStaleDays is the threshold in days before a map
	// refresh is considered stale.
	MapStaleDays = 30
	// MapStalenessThrottleID is the state file name for daily
	// throttle of map staleness checks.
	MapStalenessThrottleID = "check-map-staleness"
)
