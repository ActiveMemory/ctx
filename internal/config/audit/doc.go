//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package audit declares filesystem layout, frontmatter shape,
// status enums, and staleness thresholds for the out-of-band
// audit channel.
//
// The audit channel decouples discipline enforcement from the
// in-band commit cadence: an out-of-band auditor (a separate
// Claude Code session) drops a structured report under
// `.context/<DirName>/<kind><ReportExt>`, and the
// `ctx system checkaudit` UserPromptSubmit hook
// verbatim-relays the report's body on the next interactive
// turn. Dismissal state lives in
// `.context/<DirName>/<DismissedLedger>` so a user nuking
// `.context/state/` does not silently re-surface dismissed
// audits.
//
// See specs/audit-channel.md for the design rationale.
package audit
