//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/ActiveMemory/ctx/internal/config/loop"
	"github.com/spf13/cobra"
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

	short, long := desc.CommandDesc(cmd.DescKeyLoop)
	cmd := &cobra.Command{
		Use:   cmd.DescKeyLoop,
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
		loop.PromptMd, desc.FlagDesc(flag.DescKeyLoopPrompt),
	)
	cmd.Flags().StringVarP(
		&tool, "tool", "t", "claude", desc.FlagDesc(flag.DescKeyLoopTool),
	)
	cmd.Flags().IntVarP(
		&maxIterations,
		"max-iterations", "n",
		0, desc.FlagDesc(flag.DescKeyLoopMaxIterations),
	)
	cmd.Flags().StringVarP(
		&completionMsg,
		"completion", "c", loop.DefaultCompletionSignal,
		desc.FlagDesc(flag.DescKeyLoopCompletion),
	)
	cmd.Flags().StringVarP(
		&outputFile,
		"output", "o",
		"loop.sh", desc.FlagDesc(flag.DescKeyLoopOutput),
	)

	return cmd
}
