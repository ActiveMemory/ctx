//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stats

import "github.com/spf13/cobra"

// Table prints stats table lines. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - lines: pre-formatted stats lines (header, separator, data rows)
func Table(cmd *cobra.Command, lines []string) {
	if cmd == nil {
		return
	}
	for _, line := range lines {
		cmd.Println(line)
	}
}
