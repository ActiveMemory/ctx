//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package emit produces the shell-specific strings used by
// `ctx activate` and `ctx deactivate` to bind or clear CTX_DIR
// for the current shell via `eval "$(ctx activate)"`.
//
// v1 supports bash, zsh, and POSIX sh. All three share identical
// `export` / `unset` syntax. Fish / nushell / powershell can be
// added later by extending [Set] and [Unset] without touching the
// call sites. That extensibility is the only reason this lives in
// its own package rather than inline in the command's Run.
//
// # Supported Shells
//
//	bash, zsh, sh: POSIX export / unset
//	fish:          deferred (see specs/explicit-context-dir.md).
//
// # Detection
//
// [DetectShell] returns the first non-empty value of, in order:
// the explicit --shell flag, the basename of $SHELL, and a bash
// fallback. Users who want deterministic output in scripts should
// pass --shell explicitly.
package emit
