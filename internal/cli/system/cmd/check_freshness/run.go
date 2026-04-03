//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_freshness

import (
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreCheck "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/drift"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/state"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/freshness"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	cfgTime "github.com/ActiveMemory/ctx/internal/config/time"
	"github.com/ActiveMemory/ctx/internal/config/warn"
	"github.com/ActiveMemory/ctx/internal/entity"
	ctxLog "github.com/ActiveMemory/ctx/internal/log/warn"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Run executes the check-freshness hook logic.
//
// Reads tracked files from .ctxrc freshness_files config. For each
// file, stats it and warns if it has not been modified within the
// freshness window (~6 months). Files that do not exist are silently
// skipped. The hook is a no-op when no files are configured.
// Throttled to once per day.
//
// Parameters:
//   - cmd: Cobra command for output
//   - stdin: standard input for hook JSON
//
// Returns:
//   - error: Always nil (hook errors are non-fatal)
func Run(cmd *cobra.Command, stdin *os.File) error {
	input, _, paused := coreCheck.Preamble(stdin)
	if paused {
		return nil
	}

	files := rc.FreshnessFiles()
	if len(files) == 0 {
		return nil
	}

	tmpDir := state.Dir()
	throttleFile := filepath.Join(tmpDir, freshness.ThrottleID)

	if coreCheck.DailyThrottled(throttleFile) {
		return nil
	}

	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
		ctxLog.Warn(warn.Getwd, cwdErr)
		return nil
	}

	now := time.Now()
	var staleEntries []entity.StaleEntry

	for _, tf := range files {
		absPath := filepath.Join(cwd, tf.Path)

		info, statErr := os.Stat(absPath)
		if statErr != nil {
			continue
		}

		age := now.Sub(info.ModTime())
		if age <= freshness.StaleThreshold {
			continue
		}

		staleEntries = append(staleEntries, entity.StaleEntry{
			Path:      tf.Path,
			Desc:      tf.Desc,
			ReviewURL: tf.ReviewURL,
			Days:      int(age.Hours() / cfgTime.HoursPerDay),
		})
	}

	if len(staleEntries) == 0 {
		return nil
	}

	staleText := drift.FormatStaleEntries(staleEntries)

	vars := map[string]any{freshness.VarStaleFiles: staleText}
	nudge.LoadAndEmit(cmd,
		hook.CheckFreshness, hook.VariantStale, vars, staleText,
		desc.Text(text.DescKeyFreshnessRelayPrefix),
		desc.Text(text.DescKeyFreshnessBoxTitle),
		desc.Text(text.DescKeyFreshnessRelayMessage),
		input.SessionID, throttleFile,
	)

	return nil
}
