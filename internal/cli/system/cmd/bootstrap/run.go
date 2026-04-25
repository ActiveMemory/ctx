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
	errCtx "github.com/ActiveMemory/ctx/internal/err/context"
	"github.com/ActiveMemory/ctx/internal/rc"
	"github.com/ActiveMemory/ctx/internal/write/bootstrap"
)

// Run executes the bootstrap command, emitting context directory info,
// rules, and next steps for the calling agent.
//
// Resolution under the explicit-context-dir model
// (spec: specs/explicit-context-dir.md):
//
//   - When --context-dir or CTX_DIR is declared, bootstrap validates
//     the directory exists and then emits its usual report.
//   - When neither is declared, bootstrap returns the tailored
//     "not declared" error with a candidate-count hint. Bootstrap
//     does NOT walk to guess; walk logic lives only in
//     `ctx activate`.
//
// Parameters:
//   - cmd: Cobra command providing flags and output streams.
//
// Returns:
//   - error: non-nil if the context directory is not declared, does
//     not exist, or JSON encoding fails.
func Run(cmd *cobra.Command) error {
	dir, err := rc.ContextDir()
	if err != nil {
		cwd, _ := os.Getwd()
		return errCtx.NotDeclared(rc.ScanCandidates(cwd))
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
