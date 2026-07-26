//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package checkknowledge

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	embedFlag "github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the "ctx system check-knowledge" subcommand.
//
// Returns:
//   - *cobra.Command: Configured check-knowledge subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySystemCheckKnowledge)

	var report bool
	c := &cobra.Command{
		Use:     cmd.UseSystemCheckKnowledge,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeySystemCheckKnowledge),
		Hidden:  true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			// --report is the on-demand skill path (no stdin, no
			// throttle, no relay); bare invocation is the hook path.
			if report {
				return RunReport(cmd)
			}
			return Run(cmd, os.Stdin)
		},
	}
	flagbind.BoolFlag(c, &report,
		cFlag.Report, embedFlag.DescKeySystemCheckKnowledgeReport)
	return c
}
