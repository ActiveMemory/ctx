//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package complete

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/spf13/cobra"
)

// InfoCompletedTask reports a task marked complete.
//
// Parameters:
//   - cmd: Cobra command for output
//   - taskText: The completed task description
func InfoCompletedTask(cmd *cobra.Command, taskText string) {
	cmd.Println(fmt.Sprintf(desc.TextDesc(text.DescKeyWriteCompletedTask), taskText))
}
