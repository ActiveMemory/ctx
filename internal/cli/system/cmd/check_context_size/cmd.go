//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_context_size

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx system check-context-size" subcommand.
//
// Returns:
//   - *cobra.Command: Configured check-context-size subcommand
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc(assets.CmdDescKeySystemCheckContextSize)

	return &cobra.Command{
		Use:    "check-context-size",
		Short:  short,
		Long:   long,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, os.Stdin)
		},
	}
}
