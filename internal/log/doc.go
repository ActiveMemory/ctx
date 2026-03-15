//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package log provides append-only JSONL event logging for hook
// diagnostics. Events are written to .context/state/events.jsonl when
// enabled via event_log: true in .ctxrc. The log format is identical
// to webhook payloads (notify.Payload): one struct, two sinks.
package log
