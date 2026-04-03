//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package event provides JSONL event logging for hook lifecycle tracking.
//
// [Append] writes timestamped entries to a rotating JSONL log file in
// the context state directory. [Query] reads entries back with optional
// filters for hook name, session ID, and count limits. Log rotation
// happens automatically when the file exceeds a size threshold.
//
// Key exports: [Append], [Query].
// Used by hook handlers and the event query CLI to persist and retrieve
// lifecycle events.
package event
