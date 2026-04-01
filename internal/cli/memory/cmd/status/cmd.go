//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package status implements the "ctx memory status" subcommand.
package status

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the memory status subcommand.
//
// Returns:
//   - *cobra.Command: command for showing memory bridge status.
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyMemoryStatus)
	return &cobra.Command{
		Use:   cmd.UseMemoryStatus,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd)
		},
	}
}
