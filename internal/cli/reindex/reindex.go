//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package reindex

import (
	"github.com/spf13/cobra"

	reindexRoot "github.com/ActiveMemory/ctx/internal/cli/reindex/cmd/root"
)

// Cmd returns the reindex convenience command.
//
// Returns:
//   - *cobra.Command: The reindex command with subcommands registered
func Cmd() *cobra.Command {
	return reindexRoot.Cmd()
}
