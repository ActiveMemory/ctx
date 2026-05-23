//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package ai

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// InfoPingOK prints the success message after a
// `ctx ai ping` against the named backend succeeds.
//
// Parameters:
//   - cmd: cobra command for output. Nil is a no-op.
//   - name: the backend type label that responded.
func InfoPingOK(cmd *cobra.Command, name string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(
		desc.Text(text.DescKeyWriteAIPingOK), name,
	))
}
