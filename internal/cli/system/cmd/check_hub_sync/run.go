//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_hub_sync

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/hubsync"
	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// Run executes the check-hub-sync hook logic.
//
// If a hub connection config exists, syncs new entries from
// the hub to .context/hub/. Throttled to once per day.
// Silent when no hub is configured or no new entries.
//
// Parameters:
//   - cmd: Cobra command for output
//   - stdin: standard input for hook JSON
//
// Returns:
//   - error: Always nil (hook errors are non-fatal)
func Run(cmd *cobra.Command, stdin *os.File) error {
	_, sessionID, ctxDir, stateDir, ok := check.FullPreamble(stdin)
	bailSilently := !ok
	if bailSilently {
		return nil
	}

	connected, connErr := hubsync.Connected(ctxDir)
	if connErr != nil {
		logWarn.Warn(warn.HubConnectedProbe, connErr)
		return nil
	}
	if !connected {
		return nil
	}

	markerPath := filepath.Join(stateDir, cfgHub.ThrottleHubSync)
	if check.DailyThrottled(markerPath) {
		return nil
	}

	msg := hubsync.Sync(sessionID)
	if msg != "" {
		writeSetup.Nudge(cmd, msg)
	}
	internalIo.TouchFile(markerPath)

	return nil
}
