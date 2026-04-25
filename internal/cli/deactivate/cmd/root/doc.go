//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the `ctx deactivate` cobra command.
//
// The command emits a shell-specific `unset CTX_DIR` statement to
// stdout, paired with `ctx activate` for symmetric shell integration.
// Like activate, deactivate is in the exempt allowlist: it does not
// require a declared context directory to run (clearing CTX_DIR when
// it is already unset is a harmless no-op).
//
// Usage:
//
//	eval "$(ctx deactivate)"
//	ctx deactivate --shell zsh
package root
