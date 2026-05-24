//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package audit holds shared rendering helpers for the
// `ctx system check-audit` hook: report-body formatting
// with optional STALE prefix and the humanized age
// string used by the prefix.
//
// Reading and dismissal lookups live in
// `internal/cli/audit/core/store`; this package focuses
// solely on the verbatim-relay body shape so the hook's
// own cmd/ directory stays free of helper functions.
package audit
