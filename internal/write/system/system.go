//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package system

import "github.com/spf13/cobra"

// Line prints a single line to cmd's output. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - s: line to print
func Line(cmd *cobra.Command, s string) {
	if cmd == nil {
		return
	}
	cmd.Println(s)
}

// Lines prints multiple lines to cmd's output. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - lines: lines to print
func Lines(cmd *cobra.Command, lines []string) {
	if cmd == nil {
		return
	}
	for _, line := range lines {
		cmd.Println(line)
	}
}

// Raw prints content without a trailing newline. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - s: content to print
func Raw(cmd *cobra.Command, s string) {
	if cmd == nil {
		return
	}
	cmd.Print(s)
}

// NudgeBlock prints a nudge box followed by an empty line.
// Empty box or nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - box: formatted nudge box content
func NudgeBlock(cmd *cobra.Command, box string) {
	if cmd == nil || box == "" {
		return
	}
	cmd.Println(box)
	cmd.Println()
}
