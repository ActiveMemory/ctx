//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package decision

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/decision/cmd/add"
	"github.com/ActiveMemory/ctx/internal/cli/decision/cmd/reindex"
	"github.com/ActiveMemory/ctx/internal/cli/parent"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the decision command with subcommands.
//
// The decision command provides utilities for managing the
// DECISIONS.md file:
//   - add: Add a new decision entry
//   - reindex: Regenerate the quick-reference index
//
// Returns:
//   - *cobra.Command: The decision command with subcommands
func Cmd() *cobra.Command {
	return parent.Cmd(cmd.DescKeyDecision, cmd.UseDecision,
		add.Cmd(),
		reindex.Cmd(),
	)
}
