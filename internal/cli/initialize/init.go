//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package initialize

import (
	"github.com/spf13/cobra"

	initRoot "github.com/ActiveMemory/ctx/internal/cli/initialize/cmd/root"
)

// Cmd returns the "ctx init" command for initializing a .context/ directory.
func Cmd() *cobra.Command {
	return initRoot.Cmd()
}
