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
)

// Cmd returns the serve command.
//
// Serves a static site locally via zensical. With no argument,
// serves the journal site at .context/journal-site. With a
// directory argument, serves that directory if it contains a
// zensical.toml.
//
// This command does NOT start a ctx Hub — for that, use
// `ctx hub start`.
//
// Returns:
//   - *cobra.Command: The serve command
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyServe)

	return &cobra.Command{
		Use:     cmd.UseServe,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyServe),
		Args:    cobra.MaximumNArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			return Run(args)
		},
	}
}
