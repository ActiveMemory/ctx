//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package bootstrap

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreBootstrap "github.com/ActiveMemory/ctx/internal/cli/system/core/bootstrap"
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
// Resolution under the cwd-anchored model
// (spec: specs/cwd-anchored-context.md): `$PWD/.context/` is the
// answer or it is an error. No env var, no walk, no candidate
// scan. The typed error from rc.ContextDir is propagated unchanged
// (errCtx.ErrNoCtxHere) so the cobra layer renders it for the
// user; the only extra wrapping here is the backup-dir fallback
// for downstream stat races we cannot otherwise distinguish.
//
// Parameters:
//   - cmd: Cobra command providing flags and output streams.
//
// Returns:
//   - error: non-nil when `$PWD/.context/` is absent / not a
//     directory / unreadable, or when JSON encoding fails.
func Run(cmd *cobra.Command) error {
	dir, err := rc.ContextDir()
	if err != nil {
		return err
	}

	if _, statErr := os.Stat(dir); os.IsNotExist(statErr) {
		return errBackup.ContextDirNotFound(dir)
	}

	quiet, _ := cmd.Flags().GetBool(cFlag.Quiet)
	if quiet {
		bootstrap.Dir(cmd, dir)
		return nil
	}

	files := coreBootstrap.ListContextFiles(dir)
	rules := coreBootstrap.ParseNumberedLines(
		desc.Text(text.DescKeyBootstrapRules),
	)
	nextSteps := coreBootstrap.ParseNumberedLines(
		desc.Text(text.DescKeyBootstrapNextSteps),
	)
	warning := coreBootstrap.PluginWarning()

	jsonFlag, _ := cmd.Flags().GetBool(cFlag.JSON)
	if jsonFlag {
		bootstrap.JSON(cmd, dir, files, rules, nextSteps, warning)
		return nil
	}

	fileList := coreBootstrap.WrapFileList(
		files,
		cfgBootstrap.FileListWidth,
		cfgBootstrap.FileListIndent,
	)
	bootstrap.Text(cmd, dir, fileList, rules, nextSteps, warning)
	return nil
}
