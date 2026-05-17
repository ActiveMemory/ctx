//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package gitmeta

import (
	cfgGitmeta "github.com/ActiveMemory/ctx/internal/config/git_meta"
)

// shortSHA truncates a full SHA to the canonical short form
// length defined by
// [github.com/ActiveMemory/ctx/internal/config/git_meta.ShortLen].
//
// Parameters:
//   - s: full or already-short SHA.
//
// Returns:
//   - string: first ShortLen bytes when longer, else the input
//     unchanged.
func shortSHA(s string) string {
	if len(s) <= cfgGitmeta.ShortLen {
		return s
	}
	return s[:cfgGitmeta.ShortLen]
}
