//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/add/core/build"
	"github.com/ActiveMemory/ctx/internal/cli/add/core/validate"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/entry"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
)

// Cmd returns the "ctx decision add" subcommand.
//
// Adds a new decision entry to DECISIONS.md with the
// required provenance, context, rationale, and consequence
// flags. Implementation lives in the shared add core; this
// noun-level constructor additionally enforces that the
// three body flags are present and non-placeholder.
//
// Returns:
//   - *cobra.Command: Configured decision add subcommand
func Cmd() *cobra.Command {
	c := build.Cmd(entry.Decision, cmd.DescKeyDecisionAdd, cmd.UseDecisionAdd)
	validate.RequireBodyFlags(
		c, cFlag.Context, cFlag.Rationale, cFlag.Consequence,
	)
	return c
}
