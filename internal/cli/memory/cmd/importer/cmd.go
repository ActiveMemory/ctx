//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package importer

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/spf13/cobra"
)

// Cmd returns the memory import subcommand.
//
// Returns:
//   - *cobra.Command: command for importing MEMORY.md entries into .context/ files.
func Cmd() *cobra.Command {
	var dryRun bool

	short, long := desc.CommandDesc(cmd.DescKeyMemoryImport)
	cmd := &cobra.Command{
		Use:   "import",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, dryRun)
		},
	}

	cmd.Flags().BoolVar(
		&dryRun, "dry-run", false,
		desc.FlagDesc(flag.DescKeyMemoryImportDryRun),
	)

	return cmd
}
