//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stat

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/write/line"
)

// Table prints stats table lines. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - lines: pre-formatted stats lines (header, separator, data rows)
func Table(cmd *cobra.Command, lines []string) {
	line.All(cmd, lines)
}
