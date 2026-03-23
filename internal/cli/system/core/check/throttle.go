//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check

import (
	"os"
	"time"
)

// DailyThrottled checks if a marker file was touched today (used to
// limit certain checks to once per day).
//
// Parameters:
//   - markerPath: Absolute path to the throttle marker file
//
// Returns:
//   - bool: True if the marker was touched today
func DailyThrottled(markerPath string) bool {
	info, statErr := os.Stat(markerPath)
	if statErr != nil {
		return false
	}
	y1, m1, d1 := info.ModTime().Date()
	y2, m2, d2 := time.Now().Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
