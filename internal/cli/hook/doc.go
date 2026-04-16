//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package hook implements the ctx hook parent command.
//
// The hook command is a pure namespace that consolidates
// user-facing hook-related commands under a single CLI
// group. It has no RunE of its own and delegates all
// work to its subcommands.
//
// # Subcommands
//
//   - event: query the hook event log for past firings
//   - message: inject messages into the AI session via
//     the hook message protocol
//   - notify: send webhook notifications, set up webhook
//     URLs, and test connectivity
//   - pause: suppress all context hooks for the current
//     session
//   - resume: re-enable context hooks after a pause
//
// Hook plumbing commands (check-*, block-*, heartbeat)
// live under [internal/cli/system] rather than here,
// because they are hidden and not intended for direct
// user invocation.
//
// [Cmd] builds the parent cobra command and registers each
// subcommand (event, message, notify, pause, resume).
package hook
