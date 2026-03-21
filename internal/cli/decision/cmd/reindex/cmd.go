//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package reindex provides the "ctx decisions reindex" subcommand.
package reindex

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
)

// Cmd returns the reindex subcommand for decisions.
//
// Returns:
//   - *cobra.Command: Command for regenerating the DECISIONS.md index
func Cmd() *cobra.Command {
	short, long := desc.Command(cmd.DescKeyDecisionReindex)
	return &cobra.Command{
		Use:   cmd.UseReindex,
		Short: short,
		Long:  long,
		RunE:  Run,
	}
}
