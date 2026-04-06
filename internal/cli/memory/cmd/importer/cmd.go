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
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the memory import subcommand.
//
// Returns:
//   - *cobra.Command: command for importing MEMORY.md
//     entries into .context/ files.
func Cmd() *cobra.Command {
	var dryRun bool

	short, long := desc.Command(cmd.DescKeyMemoryImport)
	c := &cobra.Command{
		Use:     cmd.UseMemoryImport,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyMemoryImport),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, dryRun)
		},
	}

	flagbind.BoolFlag(
		c, &dryRun,
		cFlag.DryRun, flag.DescKeyMemoryImportDryRun,
	)

	return c
}
