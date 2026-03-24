//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	cfgLoop "github.com/ActiveMemory/ctx/internal/config/loop"
	"github.com/ActiveMemory/ctx/internal/err/config"
	errFs "github.com/ActiveMemory/ctx/internal/err/fs"
	"github.com/ActiveMemory/ctx/internal/write/loop"
)

// Run executes the loop command logic.
//
// Validates the tool selection, generates the loop script, and writes it
// to the output file. Prints usage instructions after generation.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - promptFile: Path to the prompt file for the AI
//   - tool: AI tool to use (claude, aider, or generic)
//   - maxIterations: Maximum loop iterations (0 for unlimited)
//   - completionMsg: Signal string that indicates loop completion
//   - outputFile: Path for the generated script
//
// Returns:
//   - error: Non-nil if the tool is invalid or file write fails
func Run(
	cmd *cobra.Command,
	promptFile, tool string,
	maxIterations int,
	completionMsg, outputFile string,
) error {
	if !cfgLoop.ValidTools[tool] {
		return config.InvalidTool(tool)
	}

	script := GenerateLoopScript(promptFile, tool, maxIterations, completionMsg)

	if writeErr := os.WriteFile(
		outputFile, []byte(script), fs.PermExec,
	); writeErr != nil {
		return errFs.FileWrite(outputFile, writeErr)
	}

	loop.InfoGenerated(
		cmd, outputFile, desc.Text(text.DescKeyHeadingLoopStart),
		tool, promptFile, maxIterations, completionMsg,
	)

	return nil
}
