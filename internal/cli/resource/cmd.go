//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resource

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// Cmd returns the "ctx resource" top-level command.
//
// Returns:
//   - *cobra.Command: Configured resource command
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyResource)

	c := &cobra.Command{
		Use:     cmd.UseResource,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyResource),
		RunE: func(c *cobra.Command, _ []string) error {
			return Run(c)
		},
	}
	c.Flags().Bool(cFlag.JSON, false,
		desc.Flag(flag.DescKeyResourceJson),
	)
	return c
}
