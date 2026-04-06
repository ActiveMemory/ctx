//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package provenance resolves session and git identity for hook
// relay output. The check-reminders hook uses these helpers to
// inject provenance lines into the user-facing nudge box.
//
// Key exports: [ShortSessionID], [DefaultVal].
package provenance
