//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/add/core/build"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/entry"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/validate"
)

// Cmd returns the "ctx decision add" subcommand.
//
// Adds a new decision entry to DECISIONS.md with the
// required provenance, context, rationale, and consequence
// flags. Implementation lives in the shared add core; this
// noun-level constructor installs a PreRunE that reads each
// body flag and calls [validate.RejectPlaceholder] to reject
// empty or placeholder values.
//
// Returns:
//   - *cobra.Command: Configured decision add subcommand
func Cmd() *cobra.Command {
	c := build.Cmd(entry.Decision, cmd.DescKeyDecisionAdd, cmd.UseDecisionAdd)
	c.PreRunE = func(cobraCmd *cobra.Command, _ []string) error {
		flags := cobraCmd.Flags()
		names := []string{
			cFlag.Context,
			cFlag.Rationale,
			cFlag.Consequence,
		}
		for _, name := range names {
			value, getErr := flags.GetString(name)
			if getErr != nil {
				return getErr
			}
			rejectErr := validate.RejectPlaceholder(name, value)
			if rejectErr != nil {
				return rejectErr
			}
		}
		return nil
	}
	return c
}
