//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package entry

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/token"
)

// FindInsertionPoint finds where to insert ctx content in an existing file.
//
// Parameters:
//   - content: Existing file content
//
// Returns:
//   - int: Position to insert at
func FindInsertionPoint(content string) int {
	lines := strings.Split(content, token.NewlineLF)
	pos := 0
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			pos += len(line) + 1
			continue
		}
		if strings.HasPrefix(trimmed, token.PrefixHeading) {
			level := 0
			for _, ch := range trimmed {
				if ch == '#' {
					level++
				} else {
					break
				}
			}
			if level == 1 {
				pos += len(line) + 1
				for j := i + 1; j < len(lines); j++ {
					if strings.TrimSpace(lines[j]) == "" {
						pos += len(lines[j]) + 1
					} else {
						break
					}
				}
				return pos
			}
			return 0
		}
		return 0
	}
	return 0
}
