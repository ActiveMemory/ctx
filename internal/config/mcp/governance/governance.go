//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package governance

import "time"

// Governance thresholds — tuned to match Claude Code hook intervals.
const (
	// DriftCheckInterval is the minimum time between drift reminders.
	DriftCheckInterval = 15 * time.Minute

	// PersistNudgeAfter is the tool call count after which a persist
	// reminder fires if no context writes have occurred.
	PersistNudgeAfter = 10

	// PersistNudgeRepeat is how often the persist nudge repeats after
	// the initial threshold.
	PersistNudgeRepeat = 8
)
