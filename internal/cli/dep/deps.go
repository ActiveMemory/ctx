//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dep

import (
	"github.com/spf13/cobra"

	depRoot "github.com/ActiveMemory/ctx/internal/cli/dep/cmd/root"
)

// Cmd returns the dep command.
func Cmd() *cobra.Command {
	return depRoot.Cmd()
}
