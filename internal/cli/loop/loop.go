//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package loop

import (
	"github.com/spf13/cobra"

	looproot "github.com/ActiveMemory/ctx/internal/cli/loop/cmd/root"
)

// Cmd returns the "ctx loop" command for generating Ralph loop scripts.
func Cmd() *cobra.Command {
	return looproot.Cmd()
}
