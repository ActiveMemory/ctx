//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package deactivate

import (
	"github.com/spf13/cobra"

	deactivateRoot "github.com/ActiveMemory/ctx/internal/cli/deactivate/cmd/root"
)

// Cmd returns the `ctx deactivate` command for registration on the
// root ctx command. See cmd/root for the full command definition.
//
// Returns:
//   - *cobra.Command: the deactivate command.
func Cmd() *cobra.Command {
	return deactivateRoot.Cmd()
}
