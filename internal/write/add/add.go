//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/spf13/cobra"
)

// InfoAddedTo confirms an entry was added to a context file.
//
// Parameters:
//   - cmd: Cobra command for output
//   - filename: Name of the file the entry was added to
func InfoAddedTo(cmd *cobra.Command, filename string) {
	cmd.Println(
		fmt.Sprintf(desc.Text(text.DescKeyWriteAddedTo), filename),
	)
}
