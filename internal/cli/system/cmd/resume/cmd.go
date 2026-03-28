//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resume

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// Cmd returns the "ctx system resume" plumbing command.
//
// Returns:
//   - *cobra.Command: Configured resume subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySystemResume)

	cmd := &cobra.Command{
		Use:    cmd.UseSystemResume,
		Short:  short,
		Long:   long,
		Hidden: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, os.Stdin)
		},
	}

	cmd.Flags().String(cFlag.SessionID, "",
		desc.Flag(flag.DescKeySystemResumeSessionId),
	)

	return cmd
}
