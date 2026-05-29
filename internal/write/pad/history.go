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

// NoHistory prints the message shown when `ctx pad undo` runs
// against an empty history directory.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func NoHistory(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.Text(text.DescKeyWritePadNoHistory))
}

// Restored prints confirmation that the pad was restored from
// a specific snapshot.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - slot: snapshot identifier (filename without extension).
func Restored(cmd *cobra.Command, slot string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWritePadRestored), slot,
	))
}
