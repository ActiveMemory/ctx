//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package ingest

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the `ctx kb ingest` command.
//
// Returns:
//   - *cobra.Command: configured command.
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyKBIngest)
	return &cobra.Command{
		Use:   cmd.UseKBIngest,
		Short: short,
		Long:  long,
		Args:  cobra.ArbitraryArgs,
		RunE:  Run,
	}
}
