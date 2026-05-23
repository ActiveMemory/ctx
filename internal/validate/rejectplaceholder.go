//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package validate

import (
	"strings"

	cfgValidate "github.com/ActiveMemory/ctx/internal/config/validate"
	errCli "github.com/ActiveMemory/ctx/internal/err/cli"
	"github.com/ActiveMemory/ctx/internal/i18n"
)

// RejectPlaceholder returns an error if value is empty,
// whitespace-only, or matches the closed placeholder set
// (TBD, see chat, n/a, etc.). Matching is case-insensitive
// after trimming. Substring matches are not rejected.
//
// Callers loop over their body flags themselves and call
// this per (flag, value) pair so the wiring is visible at
// the noun-level command's PreRunE.
//
// Parameters:
//   - flag: name of the flag, used in the error message
//   - value: raw flag value as received from cobra
//
// Returns:
//   - error: non-nil when value is empty or a placeholder; nil otherwise
func RejectPlaceholder(flag, value string) error {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return errCli.FlagEmpty(flag)
	}
	if _, hit := cfgValidate.Placeholders[i18n.Fold(trimmed)]; hit {
		return errCli.FlagPlaceholder(flag, value)
	}
	return nil
}
