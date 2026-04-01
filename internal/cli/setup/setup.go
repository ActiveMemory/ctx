//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package setup

import (
	"github.com/spf13/cobra"

	setupRoot "github.com/ActiveMemory/ctx/internal/cli/setup/cmd/root"
)

// Cmd returns the "ctx setup" command for generating AI tool integrations.
//
// Returns:
//   - *cobra.Command: The setup command with subcommands registered
func Cmd() *cobra.Command {
	return setupRoot.Cmd()
}
