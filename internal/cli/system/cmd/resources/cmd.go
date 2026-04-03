//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resources

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// Cmd returns the "ctx system resources" subcommand.
//
// Returns:
//   - *cobra.Command: Configured resources subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySystemResources)

	c := &cobra.Command{
		Use:   cmd.UseSystemResources,
		Short: short,
		Long:  long,
		RunE: func(c *cobra.Command, _ []string) error {
			return Run(c)
		},
	}
	c.Flags().Bool(cFlag.JSON, false,
		desc.Flag(flag.DescKeySystemResourcesJson),
	)
	return c
}
