//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package ceremony

// Ceremony configuration.
const (
	// ThrottleID is the state file name for daily throttle of ceremony checks.
	ThrottleID = "ceremony-reminded"
	// JournalLookback is the number of recent journal files to scan for ceremony usage.
	JournalLookback = 3
	// RememberCmd is the command name scanned in journals for /ctx-remember usage.
	RememberCmd = "ctx-remember"
	// WrapUpCmd is the command name scanned in journals for /ctx-wrap-up usage.
	WrapUpCmd = "ctx-wrap-up"
)
