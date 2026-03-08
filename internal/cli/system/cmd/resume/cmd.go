//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resume

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx system resume" plumbing command.
//
// Returns:
//   - *cobra.Command: Configured resume subcommand
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc(assets.CmdDescKeySystemResume)

	cmd := &cobra.Command{
		Use:    "resume",
		Short:  short,
		Long:   long,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, os.Stdin)
		},
	}

	cmd.Flags().String("session-id", "",
		assets.FlagDesc(assets.FlagDescKeySystemResumeSessionId),
	)

	return cmd
}
