//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

import (
	"github.com/spf13/cobra"

	hookroot "github.com/ActiveMemory/ctx/internal/cli/hook/cmd/root"
)

// Cmd returns the "ctx hook" command for generating AI tool integrations.
func Cmd() *cobra.Command {
	return hookroot.Cmd()
}
