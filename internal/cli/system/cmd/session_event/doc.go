//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sessionevent provides the hidden "ctx system session-event"
// CLI subcommand.
//
// It records session lifecycle events (start or end) to the event log
// and sends a notification. Requires --type ("start"|"end") and
// --caller (editor identifier such as "vscode") flags. No-op when the
// context directory is not initialized.
//
// Key exports: [Cmd], [Run].
// Called by editor integrations to signal session boundaries.
package sessionevent
