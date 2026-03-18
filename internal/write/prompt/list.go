//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prompt

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/spf13/cobra"
)

// PromptItem prints a single prompt name in the list.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - name: prompt template name.
func PromptItem(cmd *cobra.Command, name string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(assets.TextDesc(assets.TextDescKeyWritePromptItem), name))
}

// PromptNone prints the message when no prompts are found.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func PromptNone(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(assets.TextDesc(assets.TextDescKeyWritePromptNone))
}
