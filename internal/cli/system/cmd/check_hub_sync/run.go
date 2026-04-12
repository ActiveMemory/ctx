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
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	cfgHub "github.com/ActiveMemory/ctx/internal/config/hub"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
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
	if !state.Initialized() {
		return nil
	}

	_, sessionID, paused := check.Preamble(stdin)
	if paused {
		return nil
	}

	if !hubsync.Connected() {
		return nil
	}

	markerPath := filepath.Join(
		state.Dir(), cfgHub.ThrottleHubSync,
	)
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
