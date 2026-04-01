//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

import (
	"github.com/spf13/cobra"

	hookRoot "github.com/ActiveMemory/ctx/internal/cli/hook/cmd/root"
)

// Cmd returns the "ctx hook" command for generating AI tool integrations.
//
// Returns:
//   - *cobra.Command: The hook command with subcommands registered
func Cmd() *cobra.Command {
	return hookRoot.Cmd()
}
