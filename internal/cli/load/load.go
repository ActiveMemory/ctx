//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package load

import (
	"github.com/spf13/cobra"

	loadroot "github.com/ActiveMemory/ctx/internal/cli/load/cmd/root"
)

// Cmd returns the "ctx load" command for outputting assembled context.
func Cmd() *cobra.Command {
	return loadroot.Cmd()
}
