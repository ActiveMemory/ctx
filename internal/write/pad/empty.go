//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pad

import (
	"fmt"
	"os"

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
//   - path: key file path.
func KeyCreated(path string) {
	_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf(desc.Text(text.DescKeyWritePadKeyCreated), path))
}

func mergeSkipped(cmd *cobra.Command, dupes int) {
	if dupes == 1 {
		cmd.Println(desc.Text(text.DescKeyWritePadMergeSkipped1))
	} else {
		cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadMergeSkippedN), dupes))
	}
}
