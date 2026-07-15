//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package list implements `ctxctl audit list`: one row per
// report in `.context/audit/`, columns are id, status,
// commit-range, generated-at. A "(dismissed)" suffix on
// the status column carries the dismissal signal. Reports are
// listed newest-first so the most recent audit surfaces at the
// top of the table.
package list
