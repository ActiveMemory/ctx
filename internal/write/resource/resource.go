//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resource

import (
	"github.com/spf13/cobra"

	writeIO "github.com/ActiveMemory/ctx/internal/write/io"
)

// Text prints resource table lines. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - lines: pre-formatted resource lines (header, separator, rows, summary)
func Text(cmd *cobra.Command, lines []string) {
	writeIO.Lines(cmd, lines)
}
