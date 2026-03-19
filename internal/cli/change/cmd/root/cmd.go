//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/spf13/cobra"
)

// Cmd returns the change command.
//
// Returns:
//   - *cobra.Command: Configured change command with flags registered
func Cmd() *cobra.Command {
	var since string

	short, long := desc.CommandDesc(cmd.DescKeyChange)

	cmd := &cobra.Command{
		Use:   cmd.DescKeyChange,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, since)
		},
	}

	cmd.Flags().StringVar(&since, "since", "", desc.FlagDesc(flag.DescKeyChangesSince))

	return cmd
}
