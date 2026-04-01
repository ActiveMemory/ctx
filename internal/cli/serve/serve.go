//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package serve

import (
	"github.com/spf13/cobra"

	serveRoot "github.com/ActiveMemory/ctx/internal/cli/serve/cmd/root"
)

// Cmd returns the serve command.
//
// Returns:
//   - *cobra.Command: The serve command with subcommands registered
func Cmd() *cobra.Command {
	return serveRoot.Cmd()
}
