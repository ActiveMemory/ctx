//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package complete

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// Completed reports a task marked complete.
//
// Parameters:
//   - cmd: Cobra command for output
//   - taskText: The completed task description
func Completed(cmd *cobra.Command, taskText string) {
	cmd.Println(fmt.Sprintf(desc.Text(text.DescKeyWriteCompletedTask), taskText))
}
