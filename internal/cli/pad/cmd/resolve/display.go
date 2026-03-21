//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package resolve

import "github.com/ActiveMemory/ctx/internal/cli/pad/core"

// displayAll converts raw scratchpad entries to their display form.
//
// Parameters:
//   - entries: Raw entry strings from decryption
//
// Returns:
//   - []string: Human-readable display representations
func displayAll(entries []string) []string {
	out := make([]string, len(entries))
	for i, e := range entries {
		out[i] = core.DisplayEntry(e)
	}
	return out
}
