//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
	"github.com/spf13/cobra"
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

	short, long := desc.Command(cmd.DescKeyStatus)

	cmd := &cobra.Command{
		Use:   cmd.UseStatus,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, jsonOutput, verbose)
		},
	}

	cmd.Flags().BoolVar(
		&jsonOutput,
		"json", false, desc.Flag(flag.DescKeyStatusJson),
	)
	cmd.Flags().BoolVarP(
		&verbose, "verbose", "v", false,
		desc.Flag(flag.DescKeyStatusVerbose),
	)

	return cmd
}
