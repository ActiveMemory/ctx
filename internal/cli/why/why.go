//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package why

import (
	"github.com/spf13/cobra"

	whyRoot "github.com/ActiveMemory/ctx/internal/cli/why/cmd/root"
)

// Cmd returns the "ctx why" cobra command.
//
// Returns:
//   - *cobra.Command: The why command with subcommands registered
func Cmd() *cobra.Command {
	return whyRoot.Cmd()
}
