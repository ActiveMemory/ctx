//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package checkaudit implements the "ctxctl audit-relay"
// UserPromptSubmit hook. The hook reads every audit report
// under `.context/audit/`, drops `status: clean` and
// dismissed-against-current-digest reports, and emits a
// single verbatim-relay box concatenating the remaining
// bodies. Reports older than [cfgAudit.StalenessAge] are
// prefixed with a STALE marker but still relayed.
//
// The box copy is supplied by ctxctl as plain English Go
// constants; the relay envelope is built directly via the
// nudge package, bypassing ctx's hook-message templates.
//
// See specs/audit-channel.md and specs/ctxctl-bootstrap.md
// for the design rationale.
package checkaudit
