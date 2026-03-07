//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package memory implements the "ctx memory" CLI command group for
// bridging Claude Code's auto memory into .context/.
package memory

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the "ctx memory" parent command.
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc("memory")
	cmd := &cobra.Command{
		Use:   "memory",
		Short: short,
		Long:  long,
	}

	cmd.AddCommand(
		syncCmd(),
		statusCmd(),
		diffCmd(),
		importCmd(),
		publishCmd(),
		unpublishCmd(),
	)

	return cmd
}
