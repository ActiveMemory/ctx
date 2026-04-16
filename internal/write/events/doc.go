//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package events provides terminal output for the event
// log command (ctx events).
//
// # Exported Functions
//
// [JSON] prints pre-formatted JSONL event lines to
// stdout, one per line. This is used when the --json
// flag is set for machine-readable output.
//
// [Human] prints pre-formatted human-readable event
// lines to stdout. Each line includes a timestamp,
// event type, and summary formatted by the events
// core package.
//
// [Empty] prints a notice when the event log contains
// no entries matching the query.
//
// Both [JSON] and [Human] delegate to the shared
// [line.All] primitive for nil-safe iteration.
//
// # Message Categories
//
//   - Info: event lines in JSON or human format
//   - Empty: no-results notice
//
// # Nil Safety
//
// All functions treat a nil *cobra.Command as a no-op.
package events
