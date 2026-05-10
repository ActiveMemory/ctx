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

// Cmd returns the "ctx learning add" subcommand.
//
// Adds a new learning entry to LEARNINGS.md with the
// required provenance, context, lesson, and application
// flags. Implementation lives in the shared add core; this
// noun-level constructor installs a PreRunE that calls
// [validate.BodyFlags] to reject empty or placeholder values
// on the three body flags.
//
// Returns:
//   - *cobra.Command: Configured learning add subcommand
func Cmd() *cobra.Command {
	c := build.Cmd(entry.Learning, cmd.DescKeyLearningAdd, cmd.UseLearningAdd)
	c.PreRunE = func(cobraCmd *cobra.Command, _ []string) error {
		return validate.BodyFlags(
			cobraCmd,
			cFlag.Context, cFlag.Lesson, cFlag.Application,
		)
	}
	return c
}
