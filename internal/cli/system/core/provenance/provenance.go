//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package provenance

import (
	"github.com/spf13/cobra"

	cfgJournal "github.com/ActiveMemory/ctx/internal/config/journal"
	cfgSession "github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/exec/git"
	writeProv "github.com/ActiveMemory/ctx/internal/write/provenance"
)

// ShortSessionID truncates a session ID to ShortIDLen
// characters. Returns IDUnknown if empty.
//
// Parameters:
//   - id: Full session UUID
//
// Returns:
//   - string: Truncated ID or "unknown"
func ShortSessionID(id string) string {
	if id == "" {
		return cfgSession.IDUnknown
	}
	if len(id) > cfgJournal.ShortIDLen {
		return id[:cfgJournal.ShortIDLen]
	}
	return id
}

// Emit prints the session and git provenance line to stdout.
// This is unconditional — it runs before any hook logic.
//
// Parameters:
//   - cmd: Cobra command for output
//   - sessionID: Raw session UUID from hook input
func Emit(cmd *cobra.Command, sessionID string) {
	short := ShortSessionID(sessionID)
	branch := DefaultVal(git.CurrentBranch())
	commit := DefaultVal(git.ShortHead())

	writeProv.Line(cmd, short, branch, commit)
}

// DefaultVal returns val if non-empty, or IDUnknown.
//
// Parameters:
//   - val: Value to check
//
// Returns:
//   - string: Value or "unknown"
func DefaultVal(val string) string {
	if val == "" {
		return cfgSession.IDUnknown
	}
	return val
}
