//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package root

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the reindex convenience command.
//
// The reindex command regenerates the quick-reference index at the top of
// both DECISIONS.md and LEARNINGS.md in a single invocation.
//
// Returns:
//   - *cobra.Command: The reindex command
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc("reindex")

	return &cobra.Command{
		Use:   "reindex",
		Short: short,
		Long:  long,
		RunE:  Run,
	}
}
