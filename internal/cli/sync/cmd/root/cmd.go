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
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the "ctx sync" command for reconciling context with the codebase.
//
// The command scans the codebase for changes that should be reflected in
// context files, such as new directories, package manager files, and
// configuration files.
//
// Flags:
//   - --dry-run: Show what would change without modifying files
//
// Returns:
//   - *cobra.Command: Configured sync command with flags registered
func Cmd() *cobra.Command {
	var dryRun bool

	short, long := desc.Command(cmd.DescKeySync)

	c := &cobra.Command{
		Use:     cmd.UseSync,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeySync),
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, dryRun)
		},
	}

	flagbind.BoolFlag(c, &dryRun,
		cFlag.DryRun, flag.DescKeySyncDryRun,
	)

	return c
}
