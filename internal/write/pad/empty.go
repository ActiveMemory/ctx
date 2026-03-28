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

// Empty prints the message when the scratchpad has no entries.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func Empty(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWritePadEmpty))
}

// KeyCreated prints a key creation notice to stderr.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - path: key file path.
func KeyCreated(cmd *cobra.Command, path string) {
	if cmd == nil {
		return
	}
	cmd.PrintErrln(fmt.Sprintf(desc.Text(text.DescKeyWritePadKeyCreated), path))
}

// mergeSkipped prints a message indicating how many duplicate entries were skipped.
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
