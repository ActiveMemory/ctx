//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stats

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/cli/system/core/hook"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/spf13/cobra"

	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/rc"
	writeStats "github.com/ActiveMemory/ctx/internal/write/stat"
)

// Run executes the stats subcommand, reading and displaying per-session
// token usage statistics from JSONL state files. Supports filtering by
// session, limiting output count, JSON output, and live follow mode.
//
// Parameters:
//   - cmd: Cobra command for flag access and output
//
// Returns:
//   - error: Non-nil on stats directory read failure
func Run(cmd *cobra.Command) error {
	follow, _ := cmd.Flags().GetBool(cFlag.Follow)
	session, _ := cmd.Flags().GetString(cFlag.Session)
	last, _ := cmd.Flags().GetInt(cFlag.Last)
	jsonOut, _ := cmd.Flags().GetBool(cFlag.JSON)

	dir := filepath.Join(rc.ContextDir(), dir.State)

	entries, readErr := hook.ReadStatsDir(dir, session)
	if readErr != nil {
		return readErr
	}

	if !follow {
		writeStats.Table(cmd, hook.FormatDumpStats(entries, last, jsonOut))
		return nil
	}

	// Dump existing entries first, then stream.
	writeStats.Table(cmd, hook.FormatDumpStats(entries, last, jsonOut))

	return hook.StreamStats(cmd.OutOrStdout(), dir, session, jsonOut)
}
