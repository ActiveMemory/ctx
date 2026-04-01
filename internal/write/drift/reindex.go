//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package drift

import (
	"fmt"
	"io"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// IndexCleared writes a message indicating that the index was cleared
// because no entries were found.
//
// Parameters:
//   - w: output writer
//   - entryType: type of entry (e.g. "decisions", "learnings")
//
// Returns:
//   - error: non-nil if the write fails
func IndexCleared(w io.Writer, entryType string) error {
	_, printErr := fmt.Fprintf(
		w, desc.Text(text.DescKeyDriftCleared)+token.NewlineLF, entryType)
	return printErr
}

// IndexRegenerated writes a message indicating how many entries were
// found and regenerated in the index.
//
// Parameters:
//   - w: output writer
//   - count: number of entries regenerated
//
// Returns:
//   - error: non-nil if the write fails
func IndexRegenerated(w io.Writer, count int) error {
	_, printErr := fmt.Fprintf(
		w, desc.Text(text.DescKeyDriftRegenerated)+token.NewlineLF, count)
	return printErr
}
