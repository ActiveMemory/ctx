//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package show

import (
	"path/filepath"

	"github.com/spf13/cobra"

	coreShow "github.com/ActiveMemory/ctx/internal/cli/trace/core/show"
	"github.com/ActiveMemory/ctx/internal/config/dir"
	cfgTrace "github.com/ActiveMemory/ctx/internal/config/trace"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Run executes the trace command logic.
//
// If last > 0, shows context for the last N commits.
// If no args are given, defaults to showing the last DefaultLastShow commits.
// Otherwise shows context for the specific commit hash in args[0].
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - args: positional arguments (optional commit hash)
//   - last: number of recent commits to show (0 = use args or default)
//   - jsonOutput: whether to format output as JSON
//
// Returns:
//   - error: non-nil on execution failure
func Run(cmd *cobra.Command, args []string, last int, jsonOutput bool) error {
	contextDir, err := rc.RequireContextDir()
	if err != nil {
		cmd.SilenceUsage = true
		return err
	}
	traceDir := filepath.Join(contextDir, dir.Trace)

	if last > 0 {
		return coreShow.Last(cmd, last, contextDir, traceDir, jsonOutput)
	}

	if len(args) == 0 {
		return coreShow.Last(cmd,
			cfgTrace.DefaultLastShow, contextDir, traceDir, jsonOutput,
		)
	}

	return coreShow.Commit(cmd, args[0], contextDir, traceDir, jsonOutput)
}
