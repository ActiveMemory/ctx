//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pause

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx system pause" plumbing command.
//
// Returns:
//   - *cobra.Command: Configured pause subcommand
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc(assets.CmdDescKeySystemPause)

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
		assets.FlagDesc(assets.FlagDescKeySystemPauseSessionId),
	)
	return cmd
}
