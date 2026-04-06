//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_task_completion

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx system check-task-completion" subcommand.
//
// Returns:
//   - *cobra.Command: Configured check-task-completion subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySystemCheckTaskCompletion)

	return &cobra.Command{
		Use:     cmd.UseSystemCheckTaskCompletion,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeySystemCheckTaskCompletion),
		Hidden:  true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, os.Stdin)
		},
	}
}
