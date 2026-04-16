//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package dismiss removes reminders from the reminder store
// by ID or in bulk.
//
// # Selective Dismissal
//
// [Many] removes one or more reminders by their numeric IDs.
// All IDs are validated against the current store before any
// deletion begins, so a single invalid ID aborts the entire
// operation. This prevents partial dismissals that could
// confuse the user about which reminders remain.
//
// Validated IDs are collected into a removal set, then the
// store is rewritten with only the non-matching reminders.
// Each dismissed reminder produces a confirmation message
// via [internal/write/remind].
//
// # Bulk Dismissal
//
// [All] clears every active reminder. If the store is
// already empty, it prints a "no reminders" message and
// returns. Otherwise it prints a confirmation for each
// dismissed reminder and a summary count, then writes an
// empty slice to the store file.
package dismiss
