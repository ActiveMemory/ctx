//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package io

import "github.com/spf13/cobra"

// Lines prints each line to cmd's output. Nil cmd is a no-op.
//
// This is a shared primitive used by domain write packages to avoid
// duplicating the nil-guard + range loop. Domain packages should
// provide their own exported function that delegates here.
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
