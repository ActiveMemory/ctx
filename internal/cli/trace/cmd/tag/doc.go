//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package tag provides the "ctx trace tag" CLI subcommand.
//
// It attaches a free-text context note to a specific commit by writing
// an override entry to the trace directory. The --note flag is required
// and the commit ref is resolved to a full hash before recording.
//
// Key exports: [Cmd], [Run].
// Used by the trace command tree to register the tag subcommand.
package tag
