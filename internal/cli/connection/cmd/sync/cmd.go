//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	coreSync "github.com/ActiveMemory/ctx/internal/cli/connection/core/sync"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the connect sync subcommand.
//
// Returns:
//   - *cobra.Command: The sync subcommand
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyConnectionSync)

	return &cobra.Command{
		Use:     cmd.UseConnectionSync,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyConnectionSync),
		Args:    cobra.NoArgs,
		RunE: func(
			cobraCmd *cobra.Command, _ []string,
		) error {
			return coreSync.Run(cobraCmd)
		},
	}
}
