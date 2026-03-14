//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package pause implements the ctx system pause subcommand.
//
// It creates a session-scoped pause marker that suppresses nudge and
// reminder hooks while allowing security hooks to continue firing.
package pause
