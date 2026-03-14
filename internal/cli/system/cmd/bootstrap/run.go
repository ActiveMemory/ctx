//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"os"

	bootstrap2 "github.com/ActiveMemory/ctx/internal/config/bootstrap"
	"github.com/ActiveMemory/ctx/internal/write/bootstrap"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/cli/system/core"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// Run executes the bootstrap command, emitting context directory info,
// rules, and next steps for the calling agent.
//
// Parameters:
//   - cmd: Cobra command providing flags and output streams.
//
// Returns:
//   - error: non-nil if the context directory does not exist or JSON
//     encoding fails.
func Run(cmd *cobra.Command) error {
	dir := rc.ContextDir()

	if _, statErr := os.Stat(dir); os.IsNotExist(statErr) {
		return ctxerr.ContextDirNotFound(dir)
	}

	quiet, _ := cmd.Flags().GetBool("quiet")
	if quiet {
		cmd.Println(dir)
		return nil
	}

	files := core.ListContextFiles(dir)
	rules := core.ParseNumberedLines(
		assets.TextDesc(assets.TextDescKeyBootstrapRules),
	)
	nextSteps := core.ParseNumberedLines(
		assets.TextDesc(assets.TextDescKeyBootstrapNextSteps),
	)
	warning := core.PluginWarning()

	jsonFlag, _ := cmd.Flags().GetBool("json")
	if jsonFlag {
		return bootstrap.BootstrapJSON(cmd, dir, files, rules, nextSteps, warning)
	}

	fileList := core.WrapFileList(
		files, bootstrap2.BootstrapFileListWidth, bootstrap2.BootstrapFileListIndent,
	)
	bootstrap.BootstrapText(cmd, dir, fileList, rules, nextSteps, warning)
	return nil
}
