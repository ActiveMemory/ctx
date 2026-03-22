//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package agent

import (
	"github.com/spf13/cobra"
)

// Packet prints a rendered markdown context packet.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - content: Pre-rendered markdown packet string.
func Packet(cmd *cobra.Command, content string) {
	if cmd == nil {
		return
	}
	cmd.Print(content)
}
