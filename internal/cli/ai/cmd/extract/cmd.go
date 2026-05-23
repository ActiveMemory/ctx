//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package extract

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/flagbind"
)

// Cmd returns the `ctx ai extract` subcommand.
//
// Returns:
//   - *cobra.Command: configured extract subcommand.
func Cmd() *cobra.Command {
	var backendName string
	short, long := desc.Command(cmd.DescKeyAIExtract)
	c := &cobra.Command{
		Use:     cmd.UseAIExtract,
		Short:   short,
		Long:    long,
		Example: desc.Example(cmd.DescKeyAIExtract),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, backendName)
		},
	}
	flagbind.StringFlag(c, &backendName,
		cFlag.Backend, flag.DescKeyAIExtractBackend,
	)
	return c
}
