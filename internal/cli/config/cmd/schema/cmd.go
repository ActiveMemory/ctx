//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package schema

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/assets/read/schema"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	ctxErr "github.com/ActiveMemory/ctx/internal/err/config"
)

// Cmd returns the "ctx config schema" subcommand.
//
// Returns:
//   - *cobra.Command: Configured schema subcommand
func Cmd() *cobra.Command {
	short, long := desc.CommandDesc(cmd.DescKeyConfigSchema)

	return &cobra.Command{
		Use:   cmd.UseConfigSchema,
		Short: short,
		Long:  long,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			data, readErr := schema.Schema()
			if readErr != nil {
				return ctxErr.ReadEmbeddedSchema(readErr)
			}
			cmd.Print(string(data))
			return nil
		},
	}
}
