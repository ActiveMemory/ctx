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
// noun-level constructor additionally enforces that the
// three body flags are present and non-placeholder.
//
// Returns:
//   - *cobra.Command: Configured learning add subcommand
func Cmd() *cobra.Command {
	c := build.Cmd(entry.Learning, cmd.DescKeyLearningAdd, cmd.UseLearningAdd)
	if err := validate.RequireBodyFlags(
		c, cFlag.Context, cFlag.Lesson, cFlag.Application,
	); err != nil {
		// Programming error — see decision/cmd/add for rationale.
		panic(err)
	}
	return c
}
