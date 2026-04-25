//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package activate

import (
	"github.com/spf13/cobra"

	activateRoot "github.com/ActiveMemory/ctx/internal/cli/activate/cmd/root"
)

// Cmd returns the `ctx activate` command for registration on the
// root ctx command. See cmd/root for the full command definition.
//
// Returns:
//   - *cobra.Command: the activate command.
func Cmd() *cobra.Command {
	return activateRoot.Cmd()
}
