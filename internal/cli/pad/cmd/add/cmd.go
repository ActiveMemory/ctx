//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/spf13/cobra"
)

// Cmd returns the pad add subcommand.
//
// Returns:
//   - *cobra.Command: Configured add subcommand
func Cmd() *cobra.Command {
	var filePath string

	short, _ := desc.CommandDesc(cmd.DescKeyPadAdd)
	cmd := &cobra.Command{
		Use:   "add TEXT",
		Short: short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if filePath != "" {
				return runAddBlob(cmd, args[0], filePath)
			}
			return runAdd(cmd, args[0])
		},
	}

	cmd.Flags().StringVarP(&filePath,
		"file", "f", "",
		desc.FlagDesc(flag.DescKeyPadAddFile),
	)

	return cmd
}
