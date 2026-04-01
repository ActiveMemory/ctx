//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package nudge

// Check-persistence configuration.
const (
	// PersistencePrefix is the state file prefix for per-session
	// persistence nudge counters.
	PersistencePrefix = "persistence-nudge-"
	// PersistenceEarlyMin is the minimum prompt count before nudging begins.
	PersistenceEarlyMin = 11
	// PersistenceEarlyMax is the upper bound for the early nudge window.
	PersistenceEarlyMax = 25
	// PersistenceEarlyInterval is the number of prompts between nudges
	// during the early window (prompts 11-25).
	PersistenceEarlyInterval = 20
	// PersistenceLateInterval is the number of prompts between nudges
	// after the early window (prompts 25+).
	PersistenceLateInterval = 15
	// PersistenceLogFile is the log filename for persistence check events.
	PersistenceLogFile = "check-persistence.log"
	// PersistenceKeyCount is the state file key for prompt count.
	PersistenceKeyCount = "count"
	// KeyLastNudge is the state file key for last nudge prompt number.
	KeyLastNudge = "last_nudge"
	// PersistenceKeyLastMtime is the state file key for last modification time.
	PersistenceKeyLastMtime = "last_mtime"
)
