//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package dismiss implements `ctxctl audit dismiss`:
// removes one or more audit reports from the
// `ctxctl audit-relay` hook's relay queue. Dismissal
// is bound to the report digest at dismiss time, so a
// fresh audit run with new findings re-surfaces the
// report.
//
// Usage:
//
//	ctxctl audit dismiss <id> [<id>...]
//	ctxctl audit dismiss --all
package dismiss
