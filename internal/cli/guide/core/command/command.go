//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package command

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/write/guide"
)

// List prints all non-hidden subcommands from the root.
//
// Parameters:
//   - cmd: Cobra command for output stream and root traversal
//
// Returns:
//   - error: Always nil
func List(cmd *cobra.Command) error {
	root := cmd.Root()
	guide.CommandsHeader(cmd)
	for _, c := range root.Commands() {
		if c.Hidden {
			continue
		}
		guide.CommandLine(cmd, c.Name(), c.Short)
	}
	return nil
}
