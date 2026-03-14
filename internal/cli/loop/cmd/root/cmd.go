//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/ActiveMemory/ctx/internal/config/loop"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx loop" command for generating Ralph loop scripts.
//
// The command generates a shell script that runs an AI assistant in a loop
// until a completion signal is detected, enabling iterative development
// where the AI builds on previous work.
//
// Flags:
//   - --prompt, -p: Prompt file to use (default "PROMPT.md")
//   - --tool, -t: AI tool - claude, aider, or generic (default "claude")
//   - --max-iterations, -n: Maximum iterations, 0 for unlimited (default 0)
//   - --completion, -c: Completion signal to detect
//     (default "SYSTEM_CONVERGED")
//   - --output, -o: Output script filename (default "loop.sh")
//
// Returns:
//   - *cobra.Command: Configured loop command with flags registered
func Cmd() *cobra.Command {
	var (
		promptFile    string
		tool          string
		maxIterations int
		completionMsg string
		outputFile    string
	)

	short, long := assets.CommandDesc(assets.CmdDescKeyLoop)
	cmd := &cobra.Command{
		Use:   "loop",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(
				cmd, promptFile, tool, maxIterations, completionMsg, outputFile,
			)
		},
	}

	cmd.Flags().StringVarP(&promptFile,
		"prompt", "p",
		loop.PromptMd, assets.FlagDesc(assets.FlagDescKeyLoopPrompt),
	)
	cmd.Flags().StringVarP(
		&tool, "tool", "t", "claude", assets.FlagDesc(assets.FlagDescKeyLoopTool),
	)
	cmd.Flags().IntVarP(
		&maxIterations,
		"max-iterations", "n",
		0, assets.FlagDesc(assets.FlagDescKeyLoopMaxIterations),
	)
	cmd.Flags().StringVarP(
		&completionMsg,
		"completion", "c", loop.DefaultCompletionSignal,
		assets.FlagDesc(assets.FlagDescKeyLoopCompletion),
	)
	cmd.Flags().StringVarP(
		&outputFile,
		"output", "o",
		"loop.sh", assets.FlagDesc(assets.FlagDescKeyLoopOutput),
	)

	return cmd
}
