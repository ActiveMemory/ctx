//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/config/loop"
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

	short, long := desc.Command(cmd.DescKeyLoop)
	c := &cobra.Command{
		Use:   cmd.UseLoop,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(
				cmd, promptFile, tool, maxIterations, completionMsg, outputFile,
			)
		},
	}

	c.Flags().StringVarP(&promptFile,
		cFlag.Prompt, cFlag.ShortPrompt,
		loop.PromptMd, desc.Flag(flag.DescKeyLoopPrompt),
	)
	c.Flags().StringVarP(
		&tool, cFlag.Tool, cFlag.ShortTool,
		loop.DefaultTool, desc.Flag(flag.DescKeyLoopTool),
	)
	c.Flags().IntVarP(
		&maxIterations,
		cFlag.MaxIterations, cFlag.ShortMaxIterations,
		0, desc.Flag(flag.DescKeyLoopMaxIterations),
	)
	c.Flags().StringVarP(
		&completionMsg,
		cFlag.Completion, cFlag.ShortCompletion,
		loop.DefaultCompletionSignal,
		desc.Flag(flag.DescKeyLoopCompletion),
	)
	c.Flags().StringVarP(
		&outputFile,
		cFlag.Output, cFlag.ShortOutput,
		loop.DefaultOutput, desc.Flag(flag.DescKeyLoopOutput),
	)

	return c
}
