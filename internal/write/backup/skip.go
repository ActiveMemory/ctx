//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package backup

import (
	"fmt"
	"io"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// SkipEntry writes a message indicating that an optional archive
// entry was skipped because its source file does not exist.
//
// Parameters:
//   - w: output writer
//   - prefix: entry prefix label
func SkipEntry(w io.Writer, prefix string) {
	_, _ = fmt.Fprintf(
		w, desc.Text(text.DescKeyWriteBackupSkipEntry), prefix,
	)
}
