//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cli

import (
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets"
	"github.com/ActiveMemory/ctx/internal/config/embed"
)

// FlagRequired returns an error for a missing required flag.
//
// Parameters:
//   - name: the flag name
//
// Returns:
//   - error: "required flag \"<name>\" not set"
func FlagRequired(name string) error {
	return fmt.Errorf(
		assets.TextDesc(embed.TextDescKeyErrValidationFlagRequired), name,
	)
}

// ArgRequired returns an error for a missing required argument.
//
// Parameters:
//   - name: the argument name
//
// Returns:
//   - error: "<name> argument is required"
func ArgRequired(name string) error {
	return fmt.Errorf(
		assets.TextDesc(embed.TextDescKeyErrValidationArgRequired), name,
	)
}

// InvalidSelection returns an error for an invalid interactive
// selection.
//
// Parameters:
//   - input: the user's input
//   - max: the maximum valid selection number
//
// Returns:
//   - error: "invalid selection: <input> (expected 1-<max>)"
func InvalidSelection(input string, max int) error {
	return fmt.Errorf(
		assets.TextDesc(embed.TextDescKeyErrValidateInvalidSelection), input, max,
	)
}

// UnknownDocument returns an error for an unrecognized document alias.
//
// Parameters:
//   - alias: the unrecognized alias
//
// Returns:
//   - error: "unknown document <alias> (available: manifesto, about, invariants)"
func UnknownDocument(alias string) error {
	return fmt.Errorf(
		assets.TextDesc(embed.TextDescKeyErrValidateUnknownDocument), alias,
	)
}
