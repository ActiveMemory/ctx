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
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the "ctx loop" command for generating Ralph loop scripts.
//
// The command generates a shell script that runs an AI assistant in a loop
// until a completion signal is detected, enabling iterative development
// where the AI builds on previous work.
//
// Flags:
//   - --prompt, -p: Prompt file to use (default ".context/loop.md")
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
		Use:     cmd.UseLoop,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyLoop),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(
				cmd, promptFile, tool, maxIterations, completionMsg, outputFile,
			)
		},
	}

	flagbind.BindStringFlagsPDefault(c,
		[]*string{
			&promptFile, &tool,
			&completionMsg, &outputFile,
		},
		[]string{
			cFlag.Prompt, cFlag.Tool,
			cFlag.Completion, cFlag.Output,
		},
		[]string{
			cFlag.ShortPrompt, cFlag.ShortTool,
			cFlag.ShortCompletion, cFlag.ShortOutput,
		},
		[]string{
			loop.PromptMd, loop.DefaultTool,
			loop.DefaultCompletionSignal,
			loop.DefaultOutput,
		},
		[]string{
			flag.DescKeyLoopPrompt,
			flag.DescKeyLoopTool,
			flag.DescKeyLoopCompletion,
			flag.DescKeyLoopOutput,
		},
	)
	flagbind.IntFlagP(
		c, &maxIterations,
		cFlag.MaxIterations, cFlag.ShortMaxIterations,
		0, flag.DescKeyLoopMaxIterations,
	)

	return c
}
