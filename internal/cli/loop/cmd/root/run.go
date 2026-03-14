//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"os"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/fs"
	"github.com/spf13/cobra"

	ctxerr "github.com/ActiveMemory/ctx/internal/err"
	"github.com/ActiveMemory/ctx/internal/write"
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
	validTools := map[string]bool{"claude": true, "aider": true, "generic": true}
	if !validTools[tool] {
		return ctxerr.InvalidTool(tool)
	}

	script := GenerateLoopScript(promptFile, tool, maxIterations, completionMsg)

	if writeErr := os.WriteFile(
		outputFile, []byte(script), fs.PermExec,
	); writeErr != nil {
		return ctxerr.FileWrite(outputFile, writeErr)
	}

	write.InfoLoopGenerated(
		cmd, outputFile, assets.LoopHeadingStart,
		tool, promptFile, maxIterations, completionMsg,
	)

	return nil
}
