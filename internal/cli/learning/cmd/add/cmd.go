//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package add

import (
	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/cli/add/core/build"
	"github.com/ActiveMemory/ctx/internal/cli/add/core/jsonpayload"
	"github.com/ActiveMemory/ctx/internal/config/embed/cmd"
	"github.com/ActiveMemory/ctx/internal/config/entry"
	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/validate"
)

// Cmd returns the "ctx learning add" subcommand.
//
// Adds a new learning entry to LEARNINGS.md with the
// required provenance, context, lesson, and application
// flags. Implementation lives in the shared add core; this
// noun-level constructor installs a PreRunE that reads each
// body flag and calls [validate.RejectPlaceholder] to reject
// empty or placeholder values.
//
// Returns:
//   - *cobra.Command: Configured learning add subcommand
func Cmd() *cobra.Command {
	c := build.Cmd(entry.Learning, cmd.DescKeyLearningAdd, cmd.UseLearningAdd)
	c.PreRunE = func(cobraCmd *cobra.Command, _ []string) error {
		if overlayErr := jsonpayload.OverlayFlags(cobraCmd); overlayErr != nil {
			return overlayErr
		}
		flags := cobraCmd.Flags()
		names := []string{
			cFlag.Context,
			cFlag.Lesson,
			cFlag.Application,
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
