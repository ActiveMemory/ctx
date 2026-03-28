//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package guide

import (
	"github.com/spf13/cobra"

	guideRoot "github.com/ActiveMemory/ctx/internal/cli/guide/cmd/root"
)

// Cmd returns the "ctx guide" cobra command.
//
// Returns:
//   - *cobra.Command: The guide command with subcommands registered
func Cmd() *cobra.Command {
	return guideRoot.Cmd()
}
