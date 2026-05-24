//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package dismiss implements `ctx audit dismiss`:
// removes one or more audit reports from the
// `ctx system checkaudit` hook's relay queue. Dismissal
// is bound to the report digest at dismiss time, so a
// fresh audit run with new findings re-surfaces the
// report.
//
// Usage:
//
//	ctx audit dismiss <id> [<id>...]
//	ctx audit dismiss --all
package dismiss
