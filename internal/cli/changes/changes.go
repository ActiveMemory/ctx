//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package changes

import (
	"github.com/spf13/cobra"

	changesroot "github.com/ActiveMemory/ctx/internal/cli/changes/cmd/root"
)

// Cmd returns the changes command.
func Cmd() *cobra.Command {
	return changesroot.Cmd()
}
