//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package change

import (
	"github.com/spf13/cobra"

	changeroot "github.com/ActiveMemory/ctx/internal/cli/change/cmd/root"
)

// Cmd returns the change command.
func Cmd() *cobra.Command {
	return changeroot.Cmd()
}
