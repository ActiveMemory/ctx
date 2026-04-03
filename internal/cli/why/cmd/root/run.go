//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/why/core/menu"
	"github.com/ActiveMemory/ctx/internal/cli/why/core/show"
)

// Run dispatches to the interactive menu or direct document
// display.
//
// Parameters:
//   - cmd: Cobra command for output stream
//   - args: Command arguments; optional args[0] is the alias
//
// Returns:
//   - error: Non-nil if document not found or input invalid
func Run(
	cmd *cobra.Command, args []string,
) error {
	if len(args) == 1 {
		return show.Doc(cmd, args[0])
	}
	return menu.Show(cmd)
}
