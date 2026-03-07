//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/add/cmd/root"
)

// Cmd returns the "ctx add" command.
//
// Returns:
//   - *cobra.Command: Configured add command
func Cmd() *cobra.Command {
	return root.Cmd()
}
