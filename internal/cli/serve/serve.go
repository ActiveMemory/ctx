//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package serve

import (
	"github.com/spf13/cobra"

	serveroot "github.com/ActiveMemory/ctx/internal/cli/serve/cmd/root"
)

// Cmd returns the serve command.
func Cmd() *cobra.Command {
	return serveroot.Cmd()
}
