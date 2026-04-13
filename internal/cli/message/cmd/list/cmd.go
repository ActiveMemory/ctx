//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package list

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// Cmd returns the "ctx message list" subcommand.
//
// Returns:
//   - *cobra.Command: Configured list subcommand
func Cmd() *cobra.Command {
	short, _ := desc.Command(cmd.DescKeyMessageList)

	c := &cobra.Command{
		Use:     cmd.UseMessageList,
		Short:   short,
		Example: desc.Example(cmd.DescKeyMessageList),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}
	c.Flags().Bool(cFlag.JSON, false, desc.Flag(flag.DescKeyMessageJson))
	return c
}
