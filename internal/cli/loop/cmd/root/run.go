//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/loop/core/script"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgLoop "github.com/ActiveMemory/ctx/internal/config/loop"
	"github.com/ActiveMemory/ctx/internal/err/config"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	ctxIo "github.com/ActiveMemory/ctx/internal/io"
	"github.com/ActiveMemory/ctx/internal/write/loop"
)

// Run executes the loop command logic.
//
// Validates the tool selection, generates the loop script,
// and writes it to the output file. Prints usage instructions
// after generation.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - promptFile: Path to the prompt file for the AI
//   - tool: AI tool to use (claude, aider, or generic)
//   - regime: Context regime (standard, eocc, or sac)
//   - maxTurns: Per-iteration turn cap (0=tool default)
//   - maxIterations: Maximum loop iterations (0=unlimited)
//   - completionMsg: Signal string for loop completion
//   - outputFile: Path for the generated script
//
// Returns:
//   - error: Non-nil if tool/regime invalid or file write fails
func Run(
	cmd *cobra.Command,
	promptFile, tool, regime string,
	maxTurns, maxIterations int,
	completionMsg, outputFile string,
) error {
	if !cfgLoop.ValidTools[tool] {
		return config.InvalidTool(tool)
	}
	if regime == "" {
		regime = cfgLoop.DefaultRegime
	}
	if !cfgLoop.ValidRegimes[regime] {
		return config.InvalidRegime(regime)
	}
	if maxTurns < 0 {
		return config.NegativeMaxTurns()
	}

	s, genErr := script.Generate(script.Opts{
		PromptFile:    promptFile,
		Tool:          tool,
		Regime:        regime,
		MaxTurns:      maxTurns,
		MaxIterations: maxIterations,
		CompletionMsg: completionMsg,
	})
	if genErr != nil {
		return genErr
	}

	if writeErr := ctxIo.SafeWriteFile(
		outputFile, []byte(s), fs.PermExec,
	); writeErr != nil {
		return errFs.FileWrite(outputFile, writeErr)
	}

	loop.InfoGenerated(
		cmd, outputFile,
		desc.Text(text.DescKeyHeadingLoopStart),
		tool, promptFile, maxIterations, completionMsg,
	)

	return nil
}
