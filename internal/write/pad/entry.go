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

// EntryAdded prints confirmation that a pad entry was added.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - n: entry number (1-based).
func EntryAdded(cmd *cobra.Command, n int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadEntryAdded), n))
}

// EntryUpdated prints confirmation that a pad entry was updated.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - n: entry number (1-based).
func EntryUpdated(cmd *cobra.Command, n int) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWritePadEntryUpdated), n))
}
