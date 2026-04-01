//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package archive

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// Cmd returns the `task archive` subcommand.
//
// The archive command moves completed tasks (marked with [x]) from TASKS.md
// to a timestamped archive file in .context/archive/. Pending tasks ([ ])
// remain in TASKS.md.
//
// Flags:
//   - --dry-run: Preview changes without modifying files
//
// Returns:
//   - *cobra.Command: Configured archive subcommand
func Cmd() *cobra.Command {
	var dryRun bool

	short, long := desc.Command(cmd.DescKeyTaskArchive)

	c := &cobra.Command{
		Use:   cmd.UseTaskArchive,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, dryRun)
		},
	}

	c.Flags().BoolVar(
		&dryRun,
		cFlag.DryRun,
		false,
		desc.Flag(flag.DescKeyTaskArchiveDryRun),
	)

	return c
}
