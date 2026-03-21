//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stats

import (
	"path/filepath"

	"github.com/ActiveMemory/ctx/internal/config/dir"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	"github.com/ActiveMemory/ctx/internal/rc"
	writeStats "github.com/ActiveMemory/ctx/internal/write/stats"
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
	follow, _ := cmd.Flags().GetBool("follow")
	session, _ := cmd.Flags().GetString("session")
	last, _ := cmd.Flags().GetInt("last")
	jsonOut, _ := cmd.Flags().GetBool("json")

	dir := filepath.Join(rc.ContextDir(), dir.State)

	entries, readErr := core.ReadStatsDir(dir, session)
	if readErr != nil {
		return readErr
	}

	if !follow {
		writeStats.Table(cmd, core.FormatDumpStats(entries, last, jsonOut))
		return nil
	}

	// Dump existing entries first, then stream.
	writeStats.Table(cmd, core.FormatDumpStats(entries, last, jsonOut))

	return core.StreamStats(cmd.OutOrStdout(), dir, session, jsonOut)
}
