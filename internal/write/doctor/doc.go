//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package doctor provides terminal output for the health
// check command (ctx doctor).
//
// # Exported Functions
//
// [JSON] prints pre-marshaled JSON check results to
// stdout. This is used when the --json flag is set.
//
// [Report] prints a human-readable health report
// grouped by category. Categories include structure,
// quality, plugin, hooks, state, size, resources, and
// events. Each check result is printed with a status
// icon (pass/fail/warning) and a descriptive message.
// The report closes with a summary line showing total
// warning and error counts.
//
// # Types
//
// [ResultItem] holds the display data for a single
// check result: category name, status symbol, and
// human-readable message. The calling command maps
// domain check results into ResultItem values before
// passing them to Report.
//
// # Message Categories
//
//   - Info: per-check results with status icons
//   - Summary: warning and error totals
package doctor
