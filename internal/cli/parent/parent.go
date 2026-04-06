//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package parent

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
)

// Cmd creates a pure grouping command with desc-loaded Short/Long
// descriptions and the given subcommands registered.
//
// Parameters:
//   - descKey: assets key for Short and Long descriptions
//   - use: cobra Use string (command name)
//   - subs: subcommands to register via AddCommand
//
// Returns:
//   - *cobra.Command: configured parent command
func Cmd(descKey, use string, subs ...*cobra.Command) *cobra.Command {
	short, long := desc.Command(descKey)
	c := &cobra.Command{
		Use:     use,
		Short:   short,
		Long:    long,
		Example: desc.Example(descKey),
	}
	c.AddCommand(subs...)
	return c
}
