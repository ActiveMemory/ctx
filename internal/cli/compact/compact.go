//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package compact

import (
	"github.com/spf13/cobra"

	compactRoot "github.com/ActiveMemory/ctx/internal/cli/compact/cmd/root"
)

// Cmd returns the "ctx compact" command for cleaning up context files.
//
// Returns:
//   - *cobra.Command: The compact command with subcommands registered
func Cmd() *cobra.Command {
	return compactRoot.Cmd()
}
