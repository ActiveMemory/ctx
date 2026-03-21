//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package importer

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cflag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// Cmd returns the memory import subcommand.
//
// Returns:
//   - *cobra.Command: command for importing MEMORY.md entries into .context/ files.
func Cmd() *cobra.Command {
	var dryRun bool

	short, long := desc.Command(cmd.DescKeyMemoryImport)
	c := &cobra.Command{
		Use:   cmd.UseMemoryImport,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, dryRun)
		},
	}

	c.Flags().BoolVar(
		&dryRun, cflag.DryRun, false,
		desc.Flag(flag.DescKeyMemoryImportDryRun),
	)

	return c
}
