//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sysinfo

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// Cmd returns the "ctx sysinfo" top-level command.
//
// Returns:
//   - *cobra.Command: Configured sysinfo command
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeySysinfo)

	c := &cobra.Command{
		Use:     cmd.UseSysinfo,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeySysinfo),
		RunE: func(c *cobra.Command, _ []string) error {
			return Run(c)
		},
	}
	c.Flags().Bool(cFlag.JSON, false,
		desc.Flag(flag.DescKeySysinfoJson),
	)
	return c
}
