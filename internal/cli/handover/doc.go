//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package handover wires the `ctx handover` parent command
// and its `write` subcommand.
//
// `ctx handover write` is the per-session handover writer for
// the ctx knowledge-base editorial pipeline (Phase KB). It
// folds postdated closeouts into the new handover by default
// and archives the source closeouts; `--no-fold` skips the
// fold for mid-session checkpoints.
//
// See [github.com/ActiveMemory/ctx/internal/write/handover]
// for the writer implementation.
package handover
