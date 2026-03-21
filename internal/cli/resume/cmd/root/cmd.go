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
	cflag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// Cmd returns the top-level "ctx resume" command.
//
// Returns:
//   - *cobra.Command: Configured resume command
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyResume)

	c := &cobra.Command{
		Use:   cmd.UseResume,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			sessionID, _ := cmd.Flags().GetString(cflag.SessionID)
			return Run(cmd, sessionID)
		},
	}

	c.Flags().String(cflag.SessionID, "",
		desc.Flag(flag.DescKeyResumeSessionId),
	)

	return c
}
