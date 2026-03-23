//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package render

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// commitCount formats a commit count using singular/plural text templates.
//
// Parameters:
//   - n: number of commits
//
// Returns:
//   - string: localized count string (e.g., "1 commit", "5 commits")
func commitCount(n int) string {
	if n == 1 {
		return desc.Text(text.DescKeyTimeCommitCount)
	}
	return fmt.Sprintf(desc.Text(text.DescKeyTimeCommitsCount), n)
}
