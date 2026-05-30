//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package cli

import (
	"errors"
	"fmt"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
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
		desc.Text(text.DescKeyErrValidateFlagRequired), name,
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
		desc.Text(text.DescKeyErrValidateArgRequired), name,
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
		desc.Text(text.DescKeyErrValidateInvalidSelection), input, max,
	)
}

// UnknownDocument returns an error for an unrecognized document alias.
//
// Parameters:
//   - alias: the unrecognized alias
//
// Returns:
//   - error: "unknown document <alias>
//     (available: manifesto, about, invariants)"
func UnknownDocument(alias string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrValidateUnknownDocument), alias,
	)
}

// NoToolSpecified returns an error when no tool is configured.
//
// Returns:
//   - error: "no tool specified: use --tool <tool> or set the tool
//     field in .ctxrc"
func NoToolSpecified() error {
	return errors.New(
		desc.Text(text.DescKeyErrCliNoToolSpecified),
	)
}

// FlagEmpty returns an error for a required body flag that was
// empty or whitespace-only after trimming.
//
// Parameters:
//   - name: the flag name (without leading dashes)
//
// Returns:
//   - error: "flag --<name> is required and cannot be empty or
//     whitespace-only"
func FlagEmpty(name string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrValidateFlagEmpty), name,
	)
}

// FlagPlaceholder returns an error for a required body flag whose
// value matched the closed placeholder set.
//
// Parameters:
//   - name: the flag name (without leading dashes)
//   - value: the offending value, included verbatim in the message
//
// Returns:
//   - error: "flag --<name> rejected placeholder value <value>;
//     provide concrete content"
func FlagPlaceholder(name, value string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrValidateFlagPlaceholder),
		name, value,
	)
}

// UnknownSubcommand returns the terse error for an unrecognised
// `ctx system` subcommand. The rich, user-facing guidance (likely
// plugin/binary version skew) is emitted separately as a stdout
// relay box; this is the stderr companion that names the verb so
// the cause is legible in logs without the help dump.
//
// Parameters:
//   - verb: the unrecognised subcommand name
//
// Returns:
//   - error: "unknown ctx system subcommand \"<verb>\""
func UnknownSubcommand(verb string) error {
	return fmt.Errorf(
		desc.Text(text.DescKeyErrCliUnknownSubcommand), verb,
	)
}
