//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package watch

import (
	"github.com/spf13/cobra"

	watchRoot "github.com/ActiveMemory/ctx/internal/cli/watch/cmd/root"
)

// Cmd returns the watch command.
//
// Returns:
//   - *cobra.Command: The watch command with subcommands registered
func Cmd() *cobra.Command {
	return watchRoot.Cmd()
}
