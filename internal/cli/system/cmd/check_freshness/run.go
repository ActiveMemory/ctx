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

	hook2 "github.com/ActiveMemory/ctx/internal/cli/system/core/check"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/drift"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/entity"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/freshness"
	"github.com/ActiveMemory/ctx/internal/config/hook"
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
	input, _, paused := hook2.Preamble(stdin)
	if paused {
		return nil
	}

	files := rc.FreshnessFiles()
	if len(files) == 0 {
		return nil
	}

	tmpDir := core.StateDir()
	throttleFile := filepath.Join(tmpDir, freshness.ThrottleID)

	if core.DailyThrottled(throttleFile) {
		return nil
	}

	cwd, cwdErr := os.Getwd()
	if cwdErr != nil {
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
			Days:      int(age.Hours() / 24),
		})
	}

	if len(staleEntries) == 0 {
		return nil
	}

	staleText := drift.FormatStaleEntries(staleEntries)

	vars := map[string]any{freshness.VarStaleFiles: staleText}
	content := core.LoadMessage(hook.CheckFreshness, hook.VariantStale, vars, staleText)
	if content == "" {
		return nil
	}

	core.EmitNudge(cmd, content,
		desc.Text(text.DescKeyFreshnessRelayPrefix),
		desc.Text(text.DescKeyFreshnessBoxTitle),
		hook.CheckFreshness, hook.VariantStale,
		desc.Text(text.DescKeyFreshnessRelayMessage),
		input.SessionID, vars, throttleFile)

	return nil
}
