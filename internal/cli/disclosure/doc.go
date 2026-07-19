//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package disclosure assembles the `ctx disclosure` command group for
// progressive disclosure of the canonical knowledge files (see
// specs/progressive-disclosure.md).
//
// [Cmd] returns the parent command with its subcommand tree.
//
// Subcommands:
//
//   - ctx disclosure inspect <file>: read-only report of a root's
//     staged entries and current themes, which the dry-run digesting
//     pass consumes.
//   - ctx disclosure apply <file> --plan <path|->: the mover — moves a
//     digest plan's staged entries into per-theme files and folds their
//     gists into the root, under the append→verify→remove guards.
//
// The parse/inspect/apply domain logic lives in internal/disclosure;
// output rendering in internal/write/disclosure.
package disclosure
