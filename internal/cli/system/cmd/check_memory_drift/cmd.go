//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_memory_drift

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx system check-memory-drift" subcommand.
//
// Returns:
//   - *cobra.Command: Configured check-memory-drift subcommand
func Cmd() *cobra.Command {
	short, _ := assets.CommandDesc(assets.CmdDescKeySystemCheckMemoryDrift)

	return &cobra.Command{
		Use:    "check-memory-drift",
		Short:  short,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, os.Stdin)
		},
	}
}
