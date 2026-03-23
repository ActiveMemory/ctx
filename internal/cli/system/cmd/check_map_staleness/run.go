//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_map_staleness

import (
	"os"
	"path/filepath"
	"time"

	"github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/health"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/config/architecture"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/spf13/cobra"

	internalIo "github.com/ActiveMemory/ctx/internal/io"
	writeHook "github.com/ActiveMemory/ctx/internal/write/hook"
)

// Run executes the check-map-staleness hook logic.
//
// Reads hook input from stdin, checks the map-tracking.json file for
// stale architecture mapping, counts commits touching internal/ since
// the last refresh, and emits a relay nudge if the map is stale and
// there are relevant commits. Throttled to once per day.
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

	input, _, paused := check.Preamble(stdin)
	if paused {
		return nil
	}
	markerPath := filepath.Join(
		state.StateDir(), architecture.MapStalenessThrottleID,
	)
	if check.DailyThrottled(markerPath) {
		return nil
	}

	info := health.ReadMapTracking()
	if info == nil || info.OptedOut {
		return nil
	}

	lastRun, parseErr := time.Parse(cfgTime.DateFormat, info.LastRun)
	if parseErr != nil {
		return nil
	}

	if time.Since(lastRun) < time.Duration(architecture.MapStaleDays)*24*time.Hour {
		return nil
	}

	// Count commits touching internal/ since last run
	moduleCommits := health.CountModuleCommits(info.LastRun)
	if moduleCommits == 0 {
		return nil
	}

	dateStr := lastRun.Format(cfgTime.DateFormat)
	writeHook.Nudge(cmd, health.EmitMapStalenessWarning(input.SessionID, dateStr, moduleCommits))

	internalIo.TouchFile(markerPath)

	return nil
}
