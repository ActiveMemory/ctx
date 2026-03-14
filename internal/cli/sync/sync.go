//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package sync

import (
	"github.com/spf13/cobra"

	syncroot "github.com/ActiveMemory/ctx/internal/cli/sync/cmd/root"
)

// Cmd returns the "ctx sync" command for reconciling context with codebase.
func Cmd() *cobra.Command {
	return syncroot.Cmd()
}
