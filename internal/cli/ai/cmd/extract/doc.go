//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package extract implements `ctx ai extract`, the
// validation consumer chosen in DECISIONS.md
// 2026-05-22-220000. Reads free text from stdin, asks
// the configured AI backend to extract structured
// candidates (decisions, learnings, tasks, open
// questions) via OpenAI-style `response_format`
// json_object mode, and writes the response as a
// proposal under `.context/proposals/<TS>-extract.md`.
// Canonical `.context/*.md` files are never written
// directly.
package extract
