//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/spf13/cobra"
)

// ResolveSide prints a conflict side block: header and numbered entries.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - side: label ("OURS" or "THEIRS").
//   - entries: display strings for each entry.
func ResolveSide(cmd *cobra.Command, side string, entries []string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadResolveHeader), side))
	for i, entry := range entries {
		cmd.Println(
			fmt.Sprintf(desc.Text(text.DescKeyWritePadResolveEntry), i+1, entry),
		)
	}
}
