//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package shell holds shared string constants used by `ctx activate`
// and `ctx deactivate` when emitting shell-specific statements via
// `eval "$(ctx activate)"`.
//
// Keeping the literal identifiers (`bash`, `zsh`, `sh`), POSIX
// export/unset format strings, and single-quote escape sequences in
// internal/config/ satisfies the magic-string audit (non-config
// literals are convention violations) and consolidates the list of
// supported dialects in one place so adding fish / nushell /
// powershell becomes a single-file change.
package shell
