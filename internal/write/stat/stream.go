//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package stat

import (
	"fmt"
	"io"
)

// StreamLine writes a formatted stats line to the given writer.
//
// Parameters:
//   - w: output writer
//   - line: pre-formatted line content
func StreamLine(w io.Writer, line string) {
	_, _ = fmt.Fprintln(w, line)
}
