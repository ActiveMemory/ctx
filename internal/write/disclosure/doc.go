//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package disclosure renders the output of the `ctx disclosure` command
// family: a knowledge root's staged entries and current themes, as a
// human summary ([Human]) or machine JSON ([JSON]).
//
// # Domain
//
// It consumes the read-only Inspection value from internal/disclosure
// and writes to the cobra command's output stream. All labels are
// resolved through internal/assets/read/desc from
// commands/text/write.yaml — no English literals live here.
//
// # Concurrency
//
// Stateless formatting; safe for concurrent use given distinct writers.
package disclosure
