//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package format

import (
	cfgJournal "github.com/ActiveMemory/ctx/internal/config/journal"
	cfgSession "github.com/ActiveMemory/ctx/internal/config/session"
)

// truncateSessionID returns the first ShortIDLen characters
// of a session ID, or IDUnknown if empty.
//
// Parameters:
//   - id: Full session UUID
//
// Returns:
//   - string: Truncated ID or "unknown"
func truncateSessionID(id string) string {
	if id == "" {
		return cfgSession.IDUnknown
	}
	if len(id) > cfgJournal.ShortIDLen {
		return id[:cfgJournal.ShortIDLen]
	}
	return id
}

// defaultProvenance returns the value if non-empty, or IDUnknown.
//
// Parameters:
//   - val: Provenance value
//
// Returns:
//   - string: Value or "unknown"
func defaultProvenance(val string) string {
	if val == "" {
		return cfgSession.IDUnknown
	}
	return val
}
