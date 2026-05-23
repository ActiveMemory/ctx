//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package validate

import (
	"strings"

	errCli "github.com/ActiveMemory/ctx/internal/err/cli"
	"github.com/ActiveMemory/ctx/internal/i18n"
	"github.com/ActiveMemory/ctx/internal/rc"
)

// RejectPlaceholder returns an error if value is empty,
// whitespace-only, or matches the active placeholder set
// (TBD, see chat, n/a, etc.). Matching is Unicode
// case-fold-insensitive (via internal/i18n.Fold) after
// whitespace trimming. Only the entire trimmed input is
// checked — substring matches are not rejected.
//
// The active set is the merged result from [rc.Placeholders]:
// the shipped default locale (loaded from
// `internal/assets/i18n/placeholders/<locale>.yaml`) plus
// any user-supplied entries from `.ctxrc placeholders:`
// (EXTEND semantics — user list appended to defaults).
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
	set, loadErr := rc.Placeholders()
	if loadErr != nil {
		// The embedded defaults YAML failed to load.
		// That's a build-time invariant violation, not
		// user input. Fail closed so the operator notices.
		return errCli.FlagPlaceholder(flag, value)
	}
	if _, hit := set[i18n.Fold(trimmed)]; hit {
		return errCli.FlagPlaceholder(flag, value)
	}
	return nil
}
