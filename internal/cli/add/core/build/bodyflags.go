//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package build

import (
	"github.com/spf13/cobra"

	cFlag "github.com/ActiveMemory/ctx/internal/config/flag"
	"github.com/ActiveMemory/ctx/internal/validate"
	"github.com/ActiveMemory/ctx/internal/write/theme"
)

// RequiredBodyFlags rejects empty or placeholder values for the named
// body flags of an add subcommand.
//
// It is skipped when the add targets the themes region: `--section
// Themes` declares a theme from a "<name> — <gist>" spec, which has no
// context, rationale, consequence, lesson, or application to supply.
// Demanding them there would make declaring a theme impossible on the
// kinds that require a body.
//
// Parameters:
//   - cobraCmd: the command whose flags are being validated
//   - names: the body flag names required for this entry kind
//
// Returns:
//   - error: a flag validation error, or nil
func RequiredBodyFlags(cobraCmd *cobra.Command, names []string) error {
	flags := cobraCmd.Flags()

	section, secErr := flags.GetString(cFlag.Section)
	if secErr != nil {
		return secErr
	}
	if theme.IsTarget(section) {
		return nil
	}

	for _, name := range names {
		value, getErr := flags.GetString(name)
		if getErr != nil {
			return getErr
		}
		if rejectErr := validate.RejectPlaceholder(name, value); rejectErr != nil {
			return rejectErr
		}
	}
	return nil
}
