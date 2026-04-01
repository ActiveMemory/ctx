//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// mergeSkipped prints a message indicating how many duplicate
// entries were skipped.
//
// Parameters:
//   - cmd: Cobra command for output
//   - dupes: Number of duplicate entries that were skipped
func mergeSkipped(cmd *cobra.Command, dupes int) {
	if dupes == 1 {
		cmd.Println(desc.Text(text.DescKeyWritePadMergeSkipped1))
	} else {
		cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadMergeSkippedN), dupes))
	}
}
