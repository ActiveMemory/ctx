//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// Added confirms an entry was added to a context file.
//
// Parameters:
//   - cmd: Cobra command for output
//   - filename: Name of the file the entry was added to
func Added(cmd *cobra.Command, filename string) {
	cmd.Println(
		fmt.Sprintf(desc.Text(text.DescKeyWriteAddedTo), filename),
	)
}

// SpecNudge prints a tip suggesting a spec when appropriate.
//
// Parameters:
//   - cmd: Cobra command for output
func SpecNudge(cmd *cobra.Command) {
	cmd.Println(desc.Text(text.DescKeyWriteSpecNudgeTip))
}
