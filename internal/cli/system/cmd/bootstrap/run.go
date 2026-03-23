//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"os"

	bootstrap2 "github.com/ActiveMemory/ctx/internal/cli/system/core/bootstrap"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	cfgBootstrap "github.com/ActiveMemory/ctx/internal/config/bootstrap"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	errBackup "github.com/ActiveMemory/ctx/internal/err/backup"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write/bootstrap"
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
		return errBackup.ContextDirNotFound(dir)
	}

	quiet, _ := cmd.Flags().GetBool(cFlag.Quiet)
	if quiet {
		bootstrap.Dir(cmd, dir)
		return nil
	}

	files := bootstrap2.ListContextFiles(dir)
	rules := bootstrap2.ParseNumberedLines(
		desc.Text(text.DescKeyBootstrapRules),
	)
	nextSteps := bootstrap2.ParseNumberedLines(
		desc.Text(text.DescKeyBootstrapNextSteps),
	)
	warning := bootstrap2.PluginWarning()

	jsonFlag, _ := cmd.Flags().GetBool(cFlag.JSON)
	if jsonFlag {
		bootstrap.JSON(cmd, dir, files, rules, nextSteps, warning)
		return nil
	}

	fileList := bootstrap2.WrapFileList(
		files,
		cfgBootstrap.BootstrapFileListWidth,
		cfgBootstrap.BootstrapFileListIndent,
	)
	bootstrap.Text(cmd, dir, fileList, rules, nextSteps, warning)
	return nil
}
