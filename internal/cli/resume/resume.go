//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resume

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	resumeroot "github.com/ActiveMemory/ctx/internal/cli/resume/cmd/root"
)

// Cmd returns the top-level "ctx resume" command.
//
// Returns:
//   - *cobra.Command: Configured resume command
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc("resume")

	cmd := &cobra.Command{
		Use:   "resume",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			sessionID, _ := cmd.Flags().GetString("session-id")
			return resumeroot.Run(cmd, sessionID)
		},
	}
	cmd.Flags().String("session-id", "", "Session ID (overrides stdin)")
	return cmd
}
