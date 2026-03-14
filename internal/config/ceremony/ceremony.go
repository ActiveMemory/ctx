//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package ceremony

// Ceremony configuration.
const (
	// CeremonyThrottleID is the state file name for daily throttle of ceremony checks.
	CeremonyThrottleID = "ceremony-reminded"
	// CeremonyJournalLookback is the number of recent journal files to scan for ceremony usage.
	CeremonyJournalLookback = 3
	// CeremonyRememberCmd is the command name scanned in journals for /ctx-remember usage.
	CeremonyRememberCmd = "ctx-remember"
	// CeremonyWrapUpCmd is the command name scanned in journals for /ctx-wrap-up usage.
	CeremonyWrapUpCmd = "ctx-wrap-up"
)
