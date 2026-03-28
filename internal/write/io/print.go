//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package io

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
)

// CountLine prints a formatted count line if count is positive.
// Shared primitive to eliminate the repeated if-count-gt-zero-print pattern
// across write packages.
//
// Parameters:
//   - cmd: Cobra command for output (nil is a no-op)
//   - key: Text key for the format string (must contain one %d verb)
//   - count: Value to print; skipped when zero or negative
func CountLine(cmd *cobra.Command, key string, count int) {
	if cmd == nil || count <= 0 {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(key), count))
}

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
