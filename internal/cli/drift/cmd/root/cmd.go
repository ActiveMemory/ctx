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

// Cmd returns the "ctx drift" command for detecting stale context.
//
// The command checks for broken path references, staleness indicators,
// constitution violations, and missing required files.
//
// Flags:
//   - --json: Output results as JSON for machine parsing
//   - --fix: Auto-fix supported issues (staleness, missing_file)
//
// Returns:
//   - *cobra.Command: Configured drift command with flags registered
func Cmd() *cobra.Command {
	var (
		jsonOutput bool
		fix        bool
	)

	short, long := desc.CommandDesc(cmd.DescKeyDrift)
	cmd := &cobra.Command{
		Use:   cmd.DescKeyDrift,
		Short: short,
		Long:  long,
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(cmd, jsonOutput, fix)
		},
	}

	cmd.Flags().BoolVar(
		&jsonOutput, "json", false, desc.FlagDesc(flag.DescKeyDriftJson),
	)
	cmd.Flags().BoolVar(&fix,
		"fix", false, desc.FlagDesc(flag.DescKeyDriftFix),
	)

	return cmd
}
