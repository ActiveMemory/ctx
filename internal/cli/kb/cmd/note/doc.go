//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package note implements `ctx kb note "<text>"`.
//
// Appends a one-liner to .context/ingest/findings.md. Never
// writes to a topic page or evidence-index. Use this when you
// want to park a finding for the next ingest pass to absorb.
//
// # Refusal contract
//
//   - Empty text:
//     [github.com/ActiveMemory/ctx/internal/err/kb/cli.ErrNoteNoText].
//
// # Related packages
//
//   - [github.com/ActiveMemory/ctx/internal/config/kb/cli] supplies
//     the format and refusal strings.
package note
