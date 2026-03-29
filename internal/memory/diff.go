//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package memory

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/token"
)

// countLines returns the number of newline characters in data.
//
// Parameters:
//   - data: Raw byte content to scan
//
// Returns:
//   - int: Count of newline characters; zero for empty input
func countLines(data []byte) int {
	if len(data) == 0 {
		return 0
	}
	return bytes.Count(data, []byte(token.NewlineLF))
}

// simpleDiff produces a minimal unified-style diff header
// with added/removed lines.
//
// Parameters:
//   - oldPath: Label for the old file in the diff header
//   - newPath: Label for the new file in the diff header
//   - oldLines: Lines from the previous version
//   - newLines: Lines from the current version
//
// Returns:
//   - string: Formatted diff showing added and removed lines
func simpleDiff(oldPath, newPath string, oldLines, newLines []string) string {
	var buf strings.Builder
	_, _ = fmt.Fprintf(&buf, desc.Text(text.DescKeyMemoryDiffOldFormat), oldPath)
	_, _ = fmt.Fprintf(&buf, desc.Text(text.DescKeyMemoryDiffNewFormat), newPath)

	oldSet := make(map[string]bool, len(oldLines))
	for _, l := range oldLines {
		oldSet[l] = true
	}
	newSet := make(map[string]bool, len(newLines))
	for _, l := range newLines {
		newSet[l] = true
	}

	for _, l := range oldLines {
		if !newSet[l] {
			buf.WriteString("-" + l + token.NewlineLF)
		}
	}
	for _, l := range newLines {
		if !oldSet[l] {
			buf.WriteString("+" + l + token.NewlineLF)
		}
	}

	return buf.String()
}
