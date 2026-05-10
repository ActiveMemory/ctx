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

// RequireBodyFlags marks each named flag as cobra-required and
// wraps the command's PreRunE to reject placeholder values
// (TBD, see chat, n/a, etc., plus whitespace-only) on those
// flags. Existing PreRunE is preserved and runs after the
// placeholder check.
//
// Parameters:
//   - c: cobra command to mutate
//   - flags: names of body flags to require and policy-check
//
// Returns:
//   - error: non-nil if any named flag does not exist on c, in
//     which case the command is left unmodified
func RequireBodyFlags(c *cobra.Command, flags ...string) error {
	for _, name := range flags {
		if c.Flag(name) == nil {
			return errCli.FlagUnregistered(name, c.Name())
		}
	}
	for _, name := range flags {
		if markErr := c.MarkFlagRequired(name); markErr != nil {
			return errCli.MarkRequiredFailed(name, markErr)
		}
	}
	prev := c.PreRunE
	c.PreRunE = func(cmd *cobra.Command, args []string) error {
		for _, name := range flags {
			value, _ := cmd.Flags().GetString(name)
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
