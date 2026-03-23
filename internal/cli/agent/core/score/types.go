//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package score

import "github.com/ActiveMemory/ctx/internal/index"

// Entry is an entry block with a computed relevance score.
type Entry struct {
	index.EntryBlock
	Score  float64
	Tokens int
}
