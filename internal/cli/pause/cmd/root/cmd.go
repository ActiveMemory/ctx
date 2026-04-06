//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// Cmd returns the top-level "ctx pause" command.
//
// Returns:
//   - *cobra.Command: Configured pause command
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyPause)
	c := &cobra.Command{
		Use:     cmd.UsePause,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyPause),
		RunE: func(cmd *cobra.Command, _ []string) error {
			sessionID, _ := cmd.Flags().GetString(cFlag.SessionID)
			return Run(cmd, sessionID)
		},
	}
	c.Flags().String(cFlag.SessionID, "",
		desc.Flag(flag.DescKeyPauseSessionId),
	)
	return c
}
