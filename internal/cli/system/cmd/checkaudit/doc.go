//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package checkaudit implements the "ctx system check-audit"
// UserPromptSubmit hook. The hook reads every audit report
// under `.context/audit/`, drops `status: clean` and
// dismissed-against-current-digest reports, and emits a
// single verbatim-relay box concatenating the remaining
// bodies. Reports older than [cfgAudit.StalenessAge] are
// prefixed with a STALE marker but still relayed.
//
// See specs/audit-channel.md for the design rationale.
package checkaudit
