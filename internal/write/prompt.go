//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package write

import (
	"fmt"
	"github.com/ActiveMemory/ctx/internal/write/config"
	"github.com/spf13/cobra"
)

// PromptCreated prints the confirmation after creating a prompt template.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - name: prompt template name.
func PromptCreated(cmd *cobra.Command, name string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(config.TplPromptCreated, name))
}

// PromptNone prints the message when no prompts are found.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
func PromptNone(cmd *cobra.Command) {
	if cmd == nil {
		return
	}
	cmd.Println(config.TplPromptNone)
}

// PromptRemoved prints the confirmation after removing a prompt template.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - name: prompt template name.
func PromptRemoved(cmd *cobra.Command, name string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(config.TplPromptRemoved, name))
}

// PromptItem prints a single prompt name in the list.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - name: prompt template name.
func PromptItem(cmd *cobra.Command, name string) {
	if cmd == nil {
		return
	}
	cmd.Println(fmt.Sprintf(config.TplPromptItem, name))
}
