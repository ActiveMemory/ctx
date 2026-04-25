//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package activate implements the `ctx activate` command.
//
// Activate is the shell-integration entry point under the
// explicit-context-dir resolution model introduced in
// specs/explicit-context-dir.md. The command scans upward from CWD
// (or validates an explicit path argument) and emits a
// shell-specific `export CTX_DIR=...` statement to stdout, intended
// to be consumed via `eval "$(ctx activate)"`.
//
// Activate is the ONLY command in the CLI that walks the filesystem
// during resolution. Every other command reads the declared
// CTX_DIR / --context-dir or calls [rc.RequireContextDir] and errors
// loudly when neither is set. Centralizing walk-up in activate keeps
// silent-inference bugs confined to a single supervised entry point.
//
// # Subpackages
//
//	cmd/root : cobra command definition and resolution logic.
//	core/emit: shell-specific emitters for bash/zsh/sh.
//
// # Behavior Summary
//
// Explicit path:   strict validation (exists, is a directory,
//
//	contains CONSTITUTION.md or TASKS.md); no --force.
//
// No args:         count-based resolution: emit when exactly one
//
//	candidate is visible; refuse on zero or many.
package activate
