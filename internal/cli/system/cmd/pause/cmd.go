//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pause

import (
	"os"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx system pause" plumbing command.
//
// Returns:
//   - *cobra.Command: Configured pause subcommand
func Cmd() *cobra.Command {
	short, long := desc.CommandDesc(cmd.DescKeySystemPause)

	cmd := &cobra.Command{
		Use:    "pause",
		Short:  short,
		Long:   long,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, os.Stdin)
		},
	}
	cmd.Flags().String("session-id", "",
		desc.FlagDesc(flag.DescKeySystemPauseSessionId),
	)
	return cmd
}
