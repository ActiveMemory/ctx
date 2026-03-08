//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the status command.
//
// Flags:
//   - --json: Output as JSON for machine parsing
//   - --verbose, -v: Include file content previews
//
// Returns:
//   - *cobra.Command: Configured status command with flags registered
func Cmd() *cobra.Command {
	var (
		jsonOutput bool
		verbose    bool
	)

	short, long := assets.CommandDesc(assets.CmdDescKeyStatus)

	cmd := &cobra.Command{
		Use:   "status",
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, jsonOutput, verbose)
		},
	}

	cmd.Flags().BoolVar(
		&jsonOutput,
		"json", false, assets.FlagDesc(assets.FlagDescKeyStatusJson),
	)
	cmd.Flags().BoolVarP(
		&verbose, "verbose", "v", false,
		assets.FlagDesc(assets.FlagDescKeyStatusVerbose),
	)

	return cmd
}
