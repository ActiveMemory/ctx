//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package why

import (
	"github.com/spf13/cobra"

	whyroot "github.com/ActiveMemory/ctx/internal/cli/why/cmd/root"
)

// Cmd returns the "ctx why" cobra command.
func Cmd() *cobra.Command {
	return whyroot.Cmd()
}
