//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package audit holds shared rendering helpers for the
// `ctxctl audit-relay` hook: report-body formatting with
// optional STALE prefix and the humanized age string used
// by the prefix. The separator and STALE-prefix copy are
// supplied by ctxctl and passed in as parameters; this
// package holds no user-facing text of its own.
//
// Reading and dismissal lookups live in
// `internal/ctxctl/cli/audit/core/store`; this package focuses
// solely on the verbatim-relay body shape so the hook's
// own cmd/ directory stays free of helper functions.
package audit
