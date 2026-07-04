//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package statusline

import "github.com/spf13/cobra"

// Line prints the assembled status line to the command's output.
//
// Parameters:
//   - cmd: Cobra command for output
//   - line: rendered status line
func Line(cmd *cobra.Command, line string) {
	cmd.Println(line)
}

// Blank prints an empty status line, blanking the display when the
// status line is disabled in .ctxrc.
//
// Parameters:
//   - cmd: Cobra command for output
func Blank(cmd *cobra.Command) {
	cmd.Println()
}
