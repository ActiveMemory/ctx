//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package status

import (
	"github.com/spf13/cobra"

	statusroot "github.com/ActiveMemory/ctx/internal/cli/status/cmd/root"
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

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show context summary with token estimate",
		Long: `Display a summary of the current .context/ directory including:
  - Number of context files
  - Estimated token count
  - Status of each file
  - Recent activity

Use --verbose to include content previews for each file.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return statusroot.Run(cmd, jsonOutput, verbose)
		},
	}

	cmd.Flags().BoolVar(
		&jsonOutput,
		"json", false, "Output as JSON",
	)
	cmd.Flags().BoolVarP(
		&verbose, "verbose", "v", false, "Include file content previews",
	)

	return cmd
}
