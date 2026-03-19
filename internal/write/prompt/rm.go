//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package prompt

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// Removed prints the confirmation after removing a prompt template.
//
// Parameters:
//   - cmd: Cobra command for output. Nil is a no-op.
//   - name: prompt template name.
func Removed(cmd *cobra.Command, name string) {
	if cmd == nil {
		return
	}
	cmd.Println(
		fmt.Sprintf(desc.TextDesc(text.DescKeyWritePromptRemoved), name),
	)
}
