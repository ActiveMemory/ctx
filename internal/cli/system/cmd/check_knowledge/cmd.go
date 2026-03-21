//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_knowledge

import (
	"os"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx system check-knowledge" subcommand.
//
// Returns:
//   - *cobra.Command: Configured check-knowledge subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySystemCheckKnowledge)

	return &cobra.Command{
		Use:    "check-knowledge",
		Short:  short,
		Long:   long,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, os.Stdin)
		},
	}
}
