//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"os"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	bootstrap2 "github.com/ActiveMemory/ctx/internal/config/bootstrap"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	ctxerr "github.com/ActiveMemory/ctx/internal/err/backup"
	"github.com/ActiveMemory/ctx/internal/write/bootstrap"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/core"
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
		desc.Text(text.DescKeyBootstrapRules),
	)
	nextSteps := core.ParseNumberedLines(
		desc.Text(text.DescKeyBootstrapNextSteps),
	)
	warning := core.PluginWarning()

	jsonFlag, _ := cmd.Flags().GetBool("json")
	if jsonFlag {
		bootstrap.JSON(cmd, dir, files, rules, nextSteps, warning)
		return nil
	}

	fileList := core.WrapFileList(
		files, bootstrap2.BootstrapFileListWidth, bootstrap2.BootstrapFileListIndent,
	)
	bootstrap.Text(cmd, dir, fileList, rules, nextSteps, warning)
	return nil
}
