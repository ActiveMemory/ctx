//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package emit

import (
	"fmt"
	"strings"

	"github.com/ActiveMemory/ctx/internal/config/shell"
)

// posixSet emits `export KEY=VALUE\n` for bash/zsh/sh. Used as the
// map value in the emitters table for all POSIX-family shells.
//
// Parameters:
//   - key:         environment variable name (already well-formed).
//   - quotedValue: value wrapped by shellQuote.
//
// Returns:
//   - string: one-line export statement with trailing newline.
func posixSet(key, quotedValue string) string {
	return fmt.Sprintf(shell.FormatPOSIXExport, key, quotedValue)
}

// posixUnset emits `unset KEY\n` for bash/zsh/sh. The value
// argument is ignored (unset has no payload) but kept in the
// signature to match the emitter type.
//
// Parameters:
//   - key: environment variable name to clear.
//   - _:   unused; kept for emitter-signature compatibility.
//
// Returns:
//   - string: one-line unset statement with trailing newline.
func posixUnset(key, _ string) string {
	return fmt.Sprintf(shell.FormatPOSIXUnset, key)
}

// shellQuote wraps s in single quotes, escaping any embedded single
// quote as close-escape-reopen (`'` followed by `\'` followed by `'`).
// The resulting string is safe to paste into any POSIX-compatible
// shell regardless of s's contents.
//
// Parameters:
//   - s: raw value (typically a filesystem path).
//
// Returns:
//   - string: single-quoted, escape-safe shell literal.
func shellQuote(s string) string {
	return shell.SingleQuote +
		strings.ReplaceAll(s, shell.SingleQuote, shell.SingleQuoteEscaped) +
		shell.SingleQuote
}
