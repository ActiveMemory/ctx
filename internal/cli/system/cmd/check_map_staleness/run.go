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

	"github.com/ActiveMemory/ctx/internal/config/architecture"
	time2 "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
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
	if !core.Initialized() {
		return nil
	}

	input, _, paused := core.HookPreamble(stdin)
	if paused {
		return nil
	}
	markerPath := filepath.Join(core.StateDir(), architecture.MapStalenessThrottleID)
	if core.IsDailyThrottled(markerPath) {
		return nil
	}

	info := core.ReadMapTracking()
	if info == nil || info.OptedOut {
		return nil
	}

	lastRun, parseErr := time.Parse(time2.DateFormat, info.LastRun)
	if parseErr != nil {
		return nil
	}

	if time.Since(lastRun) < time.Duration(architecture.MapStaleDays)*24*time.Hour {
		return nil
	}

	// Count commits touching internal/ since last run
	moduleCommits := core.CountModuleCommits(info.LastRun)
	if moduleCommits == 0 {
		return nil
	}

	dateStr := lastRun.Format(time2.DateFormat)
	core.EmitMapStalenessWarning(cmd, input.SessionID, dateStr, moduleCommits)

	core.TouchFile(markerPath)

	return nil
}
