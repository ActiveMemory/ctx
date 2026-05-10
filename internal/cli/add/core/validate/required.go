//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package validate

import (
	"strings"

	"github.com/spf13/cobra"

	cfgValidate "github.com/ActiveMemory/ctx/internal/config/validate"
	errCli "github.com/ActiveMemory/ctx/internal/err/cli"
)

// BodyFlags reads each named flag from c and returns an error if
// the value is empty, whitespace-only, or matches the closed
// placeholder set (TBD, see chat, n/a, etc.). It is a pure
// function: it does not mutate c, does not register PreRunE, and
// does not call [cobra.Command.MarkFlagRequired].
//
// Callers invoke this directly from their own PreRunE so the
// wiring is visible at the noun-level command constructor.
//
// Parameters:
//   - c: cobra command whose flags are being checked
//   - flags: names of body flags to read and policy-check
//
// Returns:
//   - error: non-nil on the first flag that fails the check;
//     nil if every flag passes
func BodyFlags(c *cobra.Command, flags ...string) error {
	for _, name := range flags {
		value, getErr := c.Flags().GetString(name)
		if getErr != nil {
			return getErr
		}
		if rejectErr := RejectPlaceholder(name, value); rejectErr != nil {
			return rejectErr
		}
	}
	return nil
}

// RejectPlaceholder returns an error if value is a placeholder
// (exact case-insensitive match against the closed set, plus
// whitespace-only). Substring matches are not rejected.
//
// Parameters:
//   - flag: name of the flag, used in the error message
//   - value: raw flag value as received from cobra
//
// Returns:
//   - error: non-nil when value is a placeholder; nil otherwise
func RejectPlaceholder(flag, value string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return errCli.FlagEmpty(flag)
	}
	if _, hit := cfgValidate.Placeholders[strings.ToLower(trimmed)]; hit {
		return errCli.FlagPlaceholder(flag, value)
	}
	return nil
}
