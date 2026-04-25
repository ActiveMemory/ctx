//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check

import (
	"os"

	"github.com/spf13/cobra"

	coreSession "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	cfgSession "github.com/ActiveMemory/ctx/internal/config/session"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
)

// PausePreamble resolves the shared preamble for the pause and resume
// hooks: gate on [state.Initialized], then resolve the session ID from
// the --session-id flag, stdin JSON, and finally the IDUnknown
// fallback. Returns ok=false when the hook should bail; a probe
// warning is logged internally when the resolver reports a non-benign
// failure, so callers only need to check ok.
//
// The regular [Preamble] can't be reused here because pause / resume
// accept --session-id on the CLI (they're also callable directly by
// the user, not just by hook JSON) and they never read the pause
// counter (being asked to pause when already paused is a no-op, not a
// gate).
//
// Parameters:
//   - cmd: Cobra command for flag access.
//   - stdin: Standard input for hook JSON fallback.
//
// Returns:
//   - string: Resolved session identifier.
//   - bool: true when the caller should proceed; false when the hook
//     should bail silently.
func PausePreamble(cmd *cobra.Command, stdin *os.File) (string, bool) {
	initialized, initErr := state.Initialized()
	if initErr != nil {
		logWarn.Warn(warn.StateInitializedProbe, initErr)
		return "", false
	}
	if !initialized {
		return "", false
	}

	sessionID, _ := cmd.Flags().GetString(cFlag.SessionID)
	if sessionID == "" {
		input := coreSession.ReadInput(stdin)
		sessionID = input.SessionID
	}
	if sessionID == "" {
		sessionID = cfgSession.IDUnknown
	}
	return sessionID, true
}
