//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package add wires the cobra "ctx convention add" subcommand.
//
// The cobra wiring delegates entirely to the shared add core
// at internal/cli/add/core/build, which constructs a fully
// configured add command bound to entry.Convention. This
// package exists only to keep the noun-first command tree
// uniform: every artifact noun owns its own cmd/add child,
// even when the parent has no other subcommands today.
package add
