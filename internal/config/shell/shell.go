//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package shell holds shared string constants used by `ctx activate`
// and `ctx deactivate` when emitting shell-specific statements via
// `eval "$(ctx activate)"`.
//
// Keeping these in internal/config/ satisfies the magic-string audit
// (any non-config magic literal is a convention violation) and
// consolidates the list of supported dialects in one place so adding
// fish / nushell / powershell is a single-file change.
package shell

// Supported shell dialect identifiers (lowercase, matches
// filepath.Base($SHELL) for Unix shells).
const (
	// Bash is the POSIX-family shell identifier for GNU Bash.
	Bash = "bash"
	// Zsh is the POSIX-family shell identifier for Z shell.
	Zsh = "zsh"
	// Sh is the POSIX-family shell identifier for /bin/sh.
	Sh = "sh"
)

// Emit formats for POSIX-compatible shells (bash/zsh/sh share one
// export/unset syntax; other shells are future work).
const (
	// FormatPOSIXExport is the format string for emitting
	// `export KEY=VALUE\n`; expects (key, quotedValue).
	FormatPOSIXExport = "export %s=%s\n"
	// FormatPOSIXUnset is the format string for emitting
	// `unset KEY\n`; expects (key).
	FormatPOSIXUnset = "unset %s\n"
	// FormatStaleReplaceComment is the format used by `ctx activate`
	// to surface a stale CTX_DIR being replaced. Expects
	// (envName, oldValue, exportLine). The leading comment hash is
	// inert in `eval` output, so it is informational only.
	FormatStaleReplaceComment = "# ctx: replacing stale %s=%s\n%s"
	// FormatAlsoVisibleAdvisory is the format used by `ctx activate`
	// to surface additional .context/ candidates further up the
	// path when more than one is visible. The innermost wins
	// (selected); each additional candidate gets one of these
	// lines written to **stderr** so it actually reaches the user
	// during the standard `eval "$(ctx activate)"` invocation
	// (`eval` captures stdout but not stderr). Expects
	// (additionalPath).
	FormatAlsoVisibleAdvisory = "ctx: also visible upward: %s\n"
	// FormatActivatedAtAdvisory is the format used by `ctx activate`
	// to surface the bound .context/ path on stderr. Always
	// printed (single-candidate too) so the user always sees what
	// just got bound, not just an empty terminal. Pairs with
	// [FormatAlsoVisibleAdvisory] when multiple candidates exist.
	// Expects (selectedPath).
	FormatActivatedAtAdvisory = "ctx: activated at: %s\n"
)

// Single-quote characters used by the POSIX-shell quoting helper.
const (
	// SingleQuote wraps values in bash/zsh single-quoted strings.
	SingleQuote = "'"
	// SingleQuoteEscaped is the canonical POSIX-shell escape for an
	// embedded single quote: `'\''` (close, escape, reopen).
	SingleQuoteEscaped = `'\''`
)
