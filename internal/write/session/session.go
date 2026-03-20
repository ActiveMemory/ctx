//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package session

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/spf13/cobra"
)

// SessionPaused prints confirmation that hooks were paused.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - sessionID: the session identifier.
func SessionPaused(cmd *cobra.Command, sessionID string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.TextDesc(text.DescKeyWritePaused), sessionID))
}

// SessionResumed prints confirmation that hooks were resumed.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - sessionID: the session identifier.
func SessionResumed(cmd *cobra.Command, sessionID string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(desc.TextDesc(text.DescKeyWriteResumed), sessionID))
}

// SessionWrappedUp prints confirmation that the wrap-up marker was written.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func SessionWrappedUp(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(desc.TextDesc(text.DescKeyMarkWrappedUpConfirmed))
}
