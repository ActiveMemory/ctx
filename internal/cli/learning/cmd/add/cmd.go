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
)

// Cmd returns the "ctx learning add" subcommand.
//
// Adds a new learning entry to LEARNINGS.md with the
// required provenance, context, lesson, and application
// flags. Implementation lives in the shared add core.
//
// Returns:
//   - *cobra.Command: Configured learning add subcommand
func Cmd() *cobra.Command {
	return build.Cmd(entry.Learning, cmd.DescKeyLearningAdd, cmd.UseLearningAdd)
}
