//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pause

import (
	"github.com/spf13/cobra"

	pauseroot "github.com/ActiveMemory/ctx/internal/cli/pause/cmd/root"
)

// Cmd returns the top-level "ctx pause" command.
func Cmd() *cobra.Command {
	return pauseroot.Cmd()
}
