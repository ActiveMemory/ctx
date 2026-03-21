//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package pause

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// Confirmed prints the pause confirmation message. Nil cmd is a no-op.
//
// Parameters:
//   - cmd: Cobra command for output
//   - sessionID: session that was paused
func Confirmed(cmd *cobra.Command, sessionID string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyPauseConfirmed), sessionID))
}
