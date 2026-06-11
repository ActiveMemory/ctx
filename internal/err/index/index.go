//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package index

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// EntriesInBlock returns an error when entry bodies are found between the
// INDEX:START and INDEX:END markers, where regenerating the index would
// delete them.
//
// Parameters:
//   - fileName: Display name of the offending file (e.g., "LEARNINGS.md")
//
// Returns:
//   - error: A refusal explaining the data-loss risk and the manual fix
func EntriesInBlock(fileName string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrIndexEntriesInBlock), fileName,
	)
}

// MalformedMarkers returns an error when the INDEX:START/INDEX:END markers are
// missing, duplicated, or out of order, where regenerating the index would
// emit a second marker.
//
// Parameters:
//   - fileName: Display name of the offending file (e.g., "LEARNINGS.md")
//
// Returns:
//   - error: A refusal explaining the marker problem and the manual fix
func MalformedMarkers(fileName string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrIndexMalformedMarkers), fileName,
	)
}
