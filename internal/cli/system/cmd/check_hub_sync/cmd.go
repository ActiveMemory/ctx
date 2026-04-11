//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package check_hub_sync

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the "ctx system check-hub-sync" subcommand.
//
// This is a hidden internal hook invoked by Claude Code
// at session start to pull new entries from a registered
// ctx Hub. It is not intended for manual use.
//
// Returns:
//   - *cobra.Command: Configured check-hub-sync subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySystemCheckHubSync)

	return &cobra.Command{
		Use:     cmd.UseSystemCheckHubSync,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeySystemCheckHubSync),
		Hidden:  true,
		RunE: func(cobraCmd *cobra.Command, _ []string) error {
			return Run(cobraCmd, os.Stdin)
		},
	}
}
