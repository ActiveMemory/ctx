//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package schema

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
	ctxerr "github.com/ActiveMemory/ctx/internal/err"
)

// Cmd returns the "ctx config schema" subcommand.
//
// Returns:
//   - *cobra.Command: Configured schema subcommand
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc(assets.CmdDescKeyConfigSchema)

	return &cobra.Command{
		Use:   "schema",
		Short: short,
		Long:  long,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			data, readErr := assets.Schema()
			if readErr != nil {
				return ctxerr.ReadEmbeddedSchema(readErr)
			}
			cmd.Print(string(data))
			return nil
		},
	}
}
