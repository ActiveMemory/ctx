//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package watch

import (
	"github.com/spf13/cobra"

	watchroot "github.com/ActiveMemory/ctx/internal/cli/watch/cmd/root"
)

// Cmd returns the watch command.
func Cmd() *cobra.Command {
	return watchroot.Cmd()
}
