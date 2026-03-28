//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

import (
	"github.com/spf13/cobra"

	syncRoot "github.com/ActiveMemory/ctx/internal/cli/sync/cmd/root"
)

// Cmd returns the "ctx sync" command for reconciling context with codebase.
//
// Returns:
//   - *cobra.Command: The sync command with subcommands registered
func Cmd() *cobra.Command {
	return syncRoot.Cmd()
}
