//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package merge

import (
	"github.com/ActiveMemory/ctx/internal/config/embed"
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the pad merge subcommand.
//
// Returns:
//   - *cobra.Command: Configured merge subcommand
func Cmd() *cobra.Command {
	var keyFile string
	var dryRun bool

	short, long := assets.CommandDesc(embed.CmdDescKeyPadMerge)
	cmd := &cobra.Command{
		Use:   "merge FILE...",
		Short: short,
		Long:  long,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, args, keyFile, dryRun)
		},
	}

	cmd.Flags().StringVarP(&keyFile, "key", "k", "",
		assets.FlagDesc(embed.FlagDescKeyPadMergeKey))
	cmd.Flags().BoolVar(&dryRun, "dry-run", false,
		assets.FlagDesc(embed.FlagDescKeyPadMergeDryRun))

	return cmd
}
