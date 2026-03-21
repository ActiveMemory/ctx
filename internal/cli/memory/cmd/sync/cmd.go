//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sync implements the "ctx memory sync" subcommand.
package sync

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// Cmd returns the memory sync subcommand.
//
// Returns:
//   - *cobra.Command: command for syncing MEMORY.md to mirror.
func Cmd() *cobra.Command {
	var dryRun bool

	short, long := desc.Command(cmd.DescKeyMemorySync)
	c := &cobra.Command{
		Use:   cmd.UseMemorySync,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, dryRun)
		},
	}

	c.Flags().BoolVar(
		&dryRun, cFlag.DryRun, false, desc.Flag(flag.DescKeyMemorySyncDryRun),
	)

	return c
}
