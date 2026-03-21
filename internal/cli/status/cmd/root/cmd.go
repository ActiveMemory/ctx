//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/embed/flag"
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

	c := &cobra.Command{
		Use:   cmd.UseStatus,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return Run(cmd, jsonOutput, verbose)
		},
	}

	c.Flags().BoolVar(
		&jsonOutput,
		"json", false, desc.Flag(flag.DescKeyStatusJson),
	)
	c.Flags().BoolVarP(
		&verbose, "verbose", "v", false,
		desc.Flag(flag.DescKeyStatusVerbose),
	)

	return c
}
