//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package deactivate implements the `ctx deactivate` command.
//
// Deactivate is the counterpart to `ctx activate` under the
// explicit-context-dir resolution model. It emits a shell-specific
// `unset CTX_DIR` statement to stdout, intended for consumption via
// `eval "$(ctx deactivate)"`.
//
// The command does not touch the filesystem and does not scan for
// candidates. CTX_DIR can always be cleared safely regardless of
// which (if any) `.context/` directories are visible.
//
// # Subpackages
//
//	cmd/root : cobra command definition and run logic.
//
// # Shell Support
//
// Deactivate shares the emit package with activate
// (internal/cli/activate/core/emit) so both commands stay in
// lockstep on supported shells. v1: bash, zsh, POSIX sh.
package deactivate
