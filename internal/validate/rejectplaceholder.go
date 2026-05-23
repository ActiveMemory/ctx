//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package validate

import (
	"strings"

	"github.com/ActiveMemory/ctx/internal/assets/read/placeholders"
	"github.com/ActiveMemory/ctx/internal/config/asset"
	errCli "github.com/ActiveMemory/ctx/internal/err/cli"
	"github.com/ActiveMemory/ctx/internal/i18n"
)

// RejectPlaceholder returns an error if value is empty,
// whitespace-only, or matches the active placeholder set
// (TBD, see chat, n/a, etc.). Matching is Unicode
// case-fold-insensitive (via internal/i18n.Fold) after
// whitespace trimming. Only the entire trimmed input is
// checked — substring matches are not rejected.
//
// The active set is the shipped default for the `en`
// locale (loaded from
// `internal/assets/i18n/placeholders/en.yaml`). A future
// commit will extend this with `.ctxrc placeholders:`
// override values (EXTEND semantics: user list appended
// to defaults).
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
	set, err := placeholders.Load(asset.LocaleEN)
	if err != nil {
		// The asset is embedded, so a load failure here
		// is a build-time invariant violation, not a
		// user-facing condition. Fail closed: reject the
		// value so the operator notices.
		return errCli.FlagPlaceholder(flag, value)
	}
	if _, hit := set[i18n.Fold(trimmed)]; hit {
		return errCli.FlagPlaceholder(flag, value)
	}
	return nil
}
