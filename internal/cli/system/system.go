//   /    Context:                     https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import (
	"github.com/spf13/cobra"
)

// Cmd returns the hidden "ctx system" parent command.
//
// All subcommands implement Claude Code hook logic as native Go binaries.
// They are not intended for direct user invocation.
//
// Returns:
//   - *cobra.Command: Hidden parent command with all hook subcommands registered
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "system",
		Short:  "Internal hook commands for Claude Code",
		Hidden: true,
	}

	cmd.AddCommand(
		checkContextSizeCmd(),
		checkPersistenceCmd(),
		checkJournalCmd(),
		blockNonPathCtxCmd(),
		postCommitCmd(),
		cleanupTmpCmd(),
	)

	return cmd
}
