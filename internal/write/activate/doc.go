//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package activate carries the write-layer helpers for the
// `ctx activate` and `ctx deactivate` commands. Both produce a
// single shell-eval line (export / unset) that callers consume
// via `eval "$(ctx activate)"`.
//
// # Why a separate write package
//
// The `cmd_print` and `cmd_fprint` audits forbid `cmd.Print*` and
// `fmt.Fprint*(<user-facing stream>, ...)` outside `internal/write/`.
// The shell-eval lines are pre-formatted by
// [internal/cli/activate/core/emit] (no template substitution at the
// write layer), so this package is intentionally tiny: a single
// helper that owns the actual stdout write.
//
// # Exported Functions
//
// [Emit] writes pre-formatted shell-eval content to the cobra
// command's stdout, no trailing newline added (the emit-layer
// helpers already include one). Both `ctx activate` and
// `ctx deactivate` Run functions call it.
package activate
