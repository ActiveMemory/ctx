//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package prune implements the ctx prune top-level command.
//
// Removes stale per-session state files under .context/state/
// that have not been touched within the configured retention
// window (default 7 days).
//
// Key exports: [Cmd], [Run].
package prune
