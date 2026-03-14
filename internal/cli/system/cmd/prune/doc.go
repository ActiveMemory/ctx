//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package prune implements the ctx system prune subcommand.
//
// It removes stale per-session state files from .context/state/ that
// exceed the configured age, while preserving global state files.
package prune
