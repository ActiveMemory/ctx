//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package audit emits user-facing strings for the
// `ctxctl audit` CLI: the "no reports" sentinel, per-row
// list items, and dismissal confirmations. The English copy
// and format strings are supplied by ctxctl and passed in as
// parameters; this writer holds no text of its own. The
// relay-box formatting used by `ctxctl audit-relay` is built
// via the nudge package, not this writer.
package audit
