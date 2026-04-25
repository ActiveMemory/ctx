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

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/health"
	"github.com/ActiveMemory/ctx/internal/config/architecture"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	internalIo "github.com/ActiveMemory/ctx/internal/io"
	logWarn "github.com/ActiveMemory/ctx/internal/log/warn"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
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
	input, _, ctxDir, stateDir, ok := check.FullPreamble(stdin)
	bailSilently := !ok
	if bailSilently {
		return nil
	}
	markerPath := filepath.Join(
		stateDir, architecture.MapStalenessThrottleID,
	)
	if check.DailyThrottled(markerPath) {
		return nil
	}

	info, readErr := health.ReadMapTracking(ctxDir)
	if readErr != nil {
		logWarn.Warn(warn.ReadMapTracking, readErr)
		return nil
	}
	if info == nil || info.OptedOut {
		return nil
	}

	lastRun, parseErr := time.Parse(cfgTime.DateFormat, info.LastRun)
	if parseErr != nil {
		return nil
	}

	staleAge := time.Duration(architecture.MapStaleDays) *
		cfgTime.HoursPerDay * time.Hour
	if time.Since(lastRun) < staleAge {
		return nil
	}

	// Count commits touching internal/ since last run
	moduleCommits := health.CountModuleCommits(info.LastRun)
	if moduleCommits == 0 {
		return nil
	}

	dateStr := lastRun.Format(cfgTime.DateFormat)
	box, emitErr := health.EmitMapStalenessWarning(
		input.SessionID, dateStr, moduleCommits,
	)
	if emitErr != nil {
		return emitErr
	}
	writeSetup.Nudge(cmd, box)

	internalIo.TouchFile(markerPath)

	return nil
}
