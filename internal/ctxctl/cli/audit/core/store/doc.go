//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package store manages on-disk persistence for the audit
// channel: reports under `.context/audit/<id>.md` and the
// dismissal ledger at `.context/audit/.dismissed.json`.
//
// Reports themselves are markdown files with a YAML
// frontmatter (parsed via the sibling parse package). The
// store layer presents them as typed [Report] values and
// exposes dismissal helpers that bind a dismissal to the
// report's content digest so a fresh audit overwrite
// re-surfaces the report.
package store
