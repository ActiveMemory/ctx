//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
)

// listCommands prints all non-hidden subcommands from the root.
//
// Parameters:
//   - cmd: Cobra command for output stream and root traversal
//
// Returns:
//   - error: Always nil
func listCommands(cmd *cobra.Command) error {
	root := cmd.Root()
	cmd.Println(desc.TextDesc(text.DescKeyGuideCommandsHead))
	cmd.Println()
	for _, c := range root.Commands() {
		if c.Hidden {
			continue
		}
		cmd.Println(fmt.Sprintf(
			desc.TextDesc(text.DescKeyGuideCommandLine), c.Name(), c.Short))
	}
	return nil
}
