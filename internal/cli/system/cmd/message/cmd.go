//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package message

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/system/cmd/message/cmd/root"
)

// Cmd returns the "ctx system message" subcommand.
//
// Returns:
//   - *cobra.Command: Configured message subcommand with sub-subcommands
func Cmd() *cobra.Command {
	return root.Cmd()
}
