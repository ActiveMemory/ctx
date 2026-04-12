//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package provenance

import (
	"github.com/spf13/cobra"

	coreSession "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	cfgJournal "github.com/ActiveMemory/ctx/internal/config/journal"
	cfgSession "github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/config/stats"
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
// When session token data is available, appends a "Context: N% free"
// suffix so the agent sees how much context window remains at the
// start of each prompt.
//
// Parameters:
//   - cmd: Cobra command for output
//   - sessionID: Raw session UUID from hook input
func Emit(cmd *cobra.Command, sessionID string) {
	short := ShortSessionID(sessionID)
	branch := DefaultVal(git.CurrentBranch())
	commit := DefaultVal(git.ShortHead())
	contextSuffix := writeProv.ContextSuffix(ContextFreePct(sessionID))

	writeProv.Line(cmd, short, branch, commit, contextSuffix)
}

// ContextFreePct returns the percentage of the model's context window
// that is still free for the given session. Returns 0 when no token
// data is available (first prompt, new session, or read error), which
// the caller can treat as "no suffix to render".
//
// Parameters:
//   - sessionID: Raw session UUID from hook input
//
// Returns:
//   - int: Percentage free (1-100), or 0 when data is unavailable
func ContextFreePct(sessionID string) int {
	info, _ := coreSession.ReadTokenInfo(sessionID)
	if info.Tokens <= 0 {
		return 0
	}
	windowSize := coreSession.EffectiveContextWindow(info.Model)
	if windowSize <= 0 {
		return 0
	}
	usedPct := info.Tokens * stats.PercentMultiplier / windowSize
	if usedPct < 0 {
		usedPct = 0
	}
	if usedPct > stats.PercentMultiplier {
		usedPct = stats.PercentMultiplier
	}
	return stats.PercentMultiplier - usedPct
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
