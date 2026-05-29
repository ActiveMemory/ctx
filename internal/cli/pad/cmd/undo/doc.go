//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package undo implements the "ctx pad undo" subcommand for
// reverting the most recent destructive scratchpad mutation.
//
// # Behavior
//
// Every destructive pad operation (add, edit, mv, rm, merge,
// normalize, resolve, tag) is preceded by a snapshot of the
// existing pad blob into `.context/scratchpad.history/`. The
// undo command picks the most recent snapshot, writes the
// current pad state as a new snapshot (so a subsequent undo
// yields a redo), and promotes the picked snapshot back to
// the live pad.
//
// When no history exists yet, the command prints a friendly
// message and exits 0. This keeps cron'd or scripted use
// quiet on fresh projects.
//
// # Flags
//
// None in Phase 1. Phase 2 adds `--list`, `--to <slot>`,
// `--prune`, and `--clear`.
//
// # Delegation
//
// Snapshot and restore are handled by
// [internal/cli/pad/core/store.Restore]; user-facing output
// goes through [internal/write/pad.NoHistory] and
// [internal/write/pad.Restored]. See
// `specs/pad-undo-snapshot.md` for the full design.
package undo
