//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package audit emits user-facing strings for the
// `ctx audit` CLI: the "no reports" sentinel, per-row list
// items, and dismissal confirmations. The relay-box
// formatting used by `ctx system checkaudit` is built via
// the nudge package, not this writer.
package audit
