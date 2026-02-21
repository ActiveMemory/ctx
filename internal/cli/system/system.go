//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"github.com/spf13/cobra"
)

// Cmd returns the "ctx system" command.
//
// When invoked without a subcommand, it displays system resource usage
// (memory, swap, disk, load) with threshold-based alerts.
//
// Hidden subcommands implement Claude Code hook logic as native Go binaries
// and are not intended for direct user invocation.
//
// Returns:
//   - *cobra.Command: Parent command with RunE for resource display and all hook subcommands
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "system",
		Short: "Show system resource usage",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runResources(cmd)
		},
	}

	cmd.Flags().Bool("json", false, "Output in JSON format")

	cmd.AddCommand(
		checkContextSizeCmd(),
		checkPersistenceCmd(),
		checkJournalCmd(),
		checkCeremoniesCmd(),
		checkVersionCmd(),
		blockNonPathCtxCmd(),
		postCommitCmd(),
		cleanupTmpCmd(),
		qaReminderCmd(),
		checkResourcesCmd(),
	)

	return cmd
}
