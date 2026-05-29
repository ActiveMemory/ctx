//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package audit holds the audit-channel error taxonomy: typed
// errors for report-file I/O, frontmatter parsing, dismissal-ledger
// I/O, and CLI input validation.
//
// These errors carry data, not prose. ctxctl (tools/ctxctl) owns
// the audit channel's user-facing English and renders these typed
// errors into it at the command edge; this package holds only
// stable diagnostic codes so a Go error always has a string.
// Callers match with errors.AsType (the wrapping types) or
// errors.Is (the sentinels).
//
// See specs/ctxctl-bootstrap.md and DECISIONS.md (2026-05-27).
package audit
