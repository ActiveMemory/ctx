//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package freshness

import "time"

const (
	// StaleThreshold is the duration after which a tracked file is
	// considered stale and should be reviewed. ~6 months.
	StaleThreshold = 182 * 24 * time.Hour

	// ThrottleID is the state file name for daily throttle of
	// freshness checks.
	ThrottleID = "freshness-checked"
)
