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

// RequireBodyFlags wraps the command's PreRunE so each named flag
// is read and rejected when its value is empty, whitespace-only,
// or matches the closed placeholder set (TBD, see chat, n/a,
// etc.). Existing PreRunE is preserved and runs after the check.
//
// The check is the single enforcement point: there is no
// [cobra.Command.MarkFlagRequired] call, so help text does not
// gain a "(required)" annotation. Cobra defaults string flags to
// the empty string, which the empty check rejects with a clear
// message — making the marker redundant and the discarded error
// it returns avoidable.
//
// Parameters:
//   - c: cobra command to mutate
//   - flags: names of body flags to read and policy-check
func RequireBodyFlags(c *cobra.Command, flags ...string) {
	prev := c.PreRunE
	c.PreRunE = func(cmd *cobra.Command, args []string) error {
		for _, name := range flags {
			value, getErr := cmd.Flags().GetString(name)
			if getErr != nil {
				return getErr
			}
			if rejectErr := RejectPlaceholder(
				name, value,
			); rejectErr != nil {
				return rejectErr
			}
		}
		if prev != nil {
			return prev(cmd, args)
		}
		return nil
	}
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
