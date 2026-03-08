//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package reindex provides the "ctx learnings reindex" subcommand.
package reindex

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets"
)

// Cmd returns the reindex subcommand for learnings.
//
// Returns:
//   - *cobra.Command: Command for regenerating the LEARNINGS.md index
func Cmd() *cobra.Command {
	short, long := assets.CommandDesc(assets.CmdDescKeyLearningsReindex)
	return &cobra.Command{
		Use:   "reindex",
		Short: short,
		Long:  long,
		RunE:  Run,
	}
}
