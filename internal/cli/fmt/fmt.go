//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package fmt

import (
	"github.com/spf13/cobra"

	fmtRoot "github.com/ActiveMemory/ctx/internal/cli/fmt/cmd/root"
)

// Cmd returns the "ctx fmt" command for formatting context files.
//
// Returns:
//   - *cobra.Command: The fmt command
func Cmd() *cobra.Command {
	return fmtRoot.Cmd()
}
