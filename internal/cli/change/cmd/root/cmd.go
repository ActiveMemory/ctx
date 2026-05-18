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
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the change command.
//
// Returns:
//   - *cobra.Command: Configured change command with flags registered
func Cmd() *cobra.Command {
	var since string

	short, long := desc.Command(cmd.DescKeyChange)

	c := &cobra.Command{
		Use:     cmd.UseChange,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyChange),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, since)
		},
	}

	flagbind.StringFlag(
		c, &since,
		cFlag.Since, flag.DescKeyChangesSince,
	)

	return c
}
